package show

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/schollz/progressbar/v3"

	"go.skymyer.dev/show-control/config"
	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/dmx/driver"
	"go.skymyer.dev/show-control/library"
	"go.skymyer.dev/show-control/utils"
)

func NewFromConfig(file string, opts ...ControllerOpt) (*Controller, error) {
	var setup = &config.Setup{}
	if err := utils.LoadFromFile(file, setup); err != nil {
		return nil, err
	}
	return New(setup, opts...)
}

type ControllerOpt func(*Controller) error

func WithTickRate(rate time.Duration) ControllerOpt {
	return func(c *Controller) error {
		c.tickRate = rate
		return nil
	}
}

func WithFixtureLibrary(file string) ControllerOpt {
	return func(c *Controller) error {
		library, err := library.FixtureLibraryFromConfig(file)
		if err != nil {
			return fmt.Errorf("fixture library: %v", err)
		}
		c.library = library
		return nil
	}
}

func New(setup *config.Setup, opts ...ControllerOpt) (*Controller, error) {

	// Default configuration
	var c = &Controller{
		tickRate:      10 * time.Millisecond,
		universes:     make(map[string]*Universe),
		fixtures:      make(map[string]*Fixture),
		devices:       make(map[string]driver.Driver),
		scenes:        make(map[string]*Scene),
		effects:       make(map[string]*EffectCollection),
		activeEffects: make(map[string]bool),
		hotkeys:       make(map[string]*Hotkey),
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("option error: %v", err)
		}
	}

	// Initialize devices
	for name, device := range setup.Devices {
		driver, err := driver.New(device.Driver, device.Device)
		if err != nil {
			return nil, fmt.Errorf("cannot initialize device %q: %v", name, err)
		}
		c.devices[name] = driver
	}

	// Initialize universes
	for name, universe := range setup.Universes {

		// Create output frame
		output := dmx.NewDMX512Frame()

		// Register this output on the device
		device, found := c.devices[universe.Output.Device]
		if !found {
			return nil, fmt.Errorf("device reference %q invalid for universe %q", universe.Output.Device, name)
		}
		device.SetUniverse(0, output)

		// Register the universe
		c.universes[name] = &Universe{
			name:   name,
			output: output,
		}
	}

	// Setup fixtures
	for name, conf := range setup.Fixtures {
		if _, ok := c.universes[conf.Universe]; !ok {
			return nil, fmt.Errorf("unknown universe assignment %q for fixture %q", conf.Universe, name)
		}
		f, err := NewFixture(c.library, name, conf, c.universes[conf.Universe].GetOutput())
		if err != nil {
			return nil, fmt.Errorf("load fixture %q: %v", name, err)
		}
		c.fixtures[name] = f
	}

	// Load scenes
	for name, specs := range setup.Scenes {
		new, err := NewScene(specs, c.fixtures, setup.Groups)
		if err != nil {
			return nil, fmt.Errorf("loading scene %q: %v", name, err)
		}
		c.scenes[name] = new
	}

	// Load effects
	for name, specs := range setup.Effects {
		new, err := LoadEffect(specs, c.fixtures, setup.Groups)
		if err != nil {
			return nil, fmt.Errorf("loading effect %q: %v", name, err)
		}
		c.effects[name] = new
		c.activeEffects[name] = false
	}

	// Setup hotkeys
	for name, hk := range setup.HotKeys {
		new, err := NewHotkey(c, hk)
		if err != nil {
			return nil, fmt.Errorf("hotkey %q: %v", name, err)
		}
		c.hotkeys[name] = new
	}

	return c, nil
}

type Controller struct {
	library   library.Fixtures
	tickRate  time.Duration
	devices   map[string]driver.Driver
	universes map[string]*Universe
	fixtures  map[string]*Fixture
	scenes    map[string]*Scene
	effects   map[string]*EffectCollection

	activeScene   string
	activeEffects map[string]bool

	hotkeys map[string]*Hotkey
}

func (c *Controller) EnableHotkeys(ctx context.Context) error {
	keysEvents, err := keyboard.GetKeys(0)
	if err != nil {
		return err
	}
	go func() {
		defer keyboard.Close()

		for {
			event := <-keysEvents
			if event.Err != nil {
				panic(event.Err)
			}
			//fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
			if event.Key == keyboard.KeyCtrlC {
				ctx.Done()
				break
			}

			if hk, ok := c.hotkeys[string(event.Rune)]; ok {
				hk.Handle()
			}
		}
	}()
	return nil
}

func NewHotkey(c *Controller, cfg *config.Hotkey) (*Hotkey, error) {
	switch cfg.Kind {
	case config.HOTKEY_KIND_SCENE:
		return &Hotkey{cfg, c.SetActiveScene}, nil
	case config.HOTKEY_KIND_EFFECT:
		return &Hotkey{cfg, c.ToggleEffect}, nil
	}
	return nil, fmt.Errorf("unknown kind %q", cfg.Kind)
}

type Hotkey struct {
	c  *config.Hotkey
	fn EventHandler
}

func (h *Hotkey) Handle() error {
	return h.fn(h.c.Value)
}

type EventHandler func(in string) error

func (c *Controller) Run(ctx context.Context) error {

	// Startup device output
	for _, device := range c.devices {
		device.Run(ctx)
	}

	// TODO: proper init & default values
	dmx.FrameLock.Lock()

	c.fixtures["sl"].MustSelectFunction("settings", "blackout-gobo-off")
	c.fixtures["sr"].MustSelectFunction("settings", "blackout-gobo-off")

	dmx.FrameLock.Unlock()

	// Live progress bar
	bar := progressbar.NewOptions(-1,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionShowIts(), // Show iterations per second
		progressbar.OptionSetItsString("renders"),
		progressbar.OptionSetDescription("[green]Running[reset]"),
		progressbar.OptionSpinnerType(23),
		progressbar.OptionThrottle(c.tickRate*10),
	)

	var (
		washDim = byte(64)

		scanMacro      = "m12" // m6
		scanMacroSpeed = byte(0)

		scanDim          = byte(255)
		scanColor        = "blue"
		scanPan          = byte(128)
		scanTilt         = byte(32)
		scanPanTiltSpeed = byte(0)

		mhAutoMode = "off"
		mhColor    = "blue"
		mhGobo     = "triangle"

		mhDim          = byte(255)
		mhPan          = byte(128)
		mhTilt         = byte(128)
		mhPanTiltSpeed = byte(0)
	)

	// Startup short handler
	c.EnableHotkeys(ctx)

	// Main show controller handler
	ticker := time.NewTicker(c.tickRate)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				for _, e := range c.effects {
					e.Stop()
				}
				bar.Close()
				return
			case <-ticker.C:
				start := time.Now()
				bar.Add(1)
				dmx.FrameLock.Lock()

				c.fixtures["laser"].MustSelectFunction("mode", "off")

				c.fixtures["moonl"].MustSelectFunction("uv", "off")
				c.fixtures["moonr"].MustSelectFunction("uv", "off")

				c.fixtures["moonl"].MustSelectFunction("beam", "off")
				c.fixtures["moonr"].MustSelectFunction("beam", "off")

				c.fixtures["moonl"].MustSetValue("speed", "speed", byte(200))
				c.fixtures["moonr"].MustSetValue("speed", "speed", byte(200))

				// cleanup
				c.fixtures["wl"].MustSetValue("dimmer", "dimmer", washDim)
				c.fixtures["wr"].MustSetValue("dimmer", "dimmer", washDim)

				c.fixtures["sl"].MustSetValue("dimmer", "dimmer", scanDim)
				c.fixtures["sr"].MustSetValue("dimmer", "dimmer", scanDim)

				c.fixtures["sl"].MustSetValue("speed", "pantilt", scanPanTiltSpeed)
				c.fixtures["sr"].MustSetValue("speed", "pantilt", scanPanTiltSpeed)

				c.fixtures["sl"].MustSetValue("move-macro-speed", "speed", scanMacroSpeed)
				c.fixtures["sr"].MustSetValue("move-macro-speed", "speed", scanMacroSpeed)

				c.fixtures["sl"].MustSelectFunction("move-macro", scanMacro)
				c.fixtures["sr"].MustSelectFunction("move-macro", scanMacro)

				c.fixtures["sl"].MustSelectFunction("gobo", scanColor)
				c.fixtures["sr"].MustSelectFunction("gobo", scanColor)

				c.fixtures["sl"].MustSetValue("pan", "pan", scanPan)
				c.fixtures["sr"].MustSetValue("pan", "pan", scanPan)

				c.fixtures["sl"].MustSetValue("tilt", "tilt", scanTilt)
				c.fixtures["sr"].MustSetValue("tilt", "tilt", scanTilt)

				c.fixtures["mh1"].MustSetValue("speed", "pantilt", mhPanTiltSpeed)
				c.fixtures["mh2"].MustSetValue("speed", "pantilt", mhPanTiltSpeed)
				c.fixtures["mh3"].MustSetValue("speed", "pantilt", mhPanTiltSpeed)
				c.fixtures["mh4"].MustSetValue("speed", "pantilt", mhPanTiltSpeed)

				c.fixtures["mh1"].MustSetValue("dimmer", "dimmer", mhDim)
				c.fixtures["mh2"].MustSetValue("dimmer", "dimmer", mhDim)
				c.fixtures["mh3"].MustSetValue("dimmer", "dimmer", mhDim)
				c.fixtures["mh4"].MustSetValue("dimmer", "dimmer", mhDim)

				c.fixtures["mh1"].MustSelectFunction("auto", mhAutoMode)
				c.fixtures["mh2"].MustSelectFunction("auto", mhAutoMode)
				c.fixtures["mh3"].MustSelectFunction("auto", mhAutoMode)
				c.fixtures["mh4"].MustSelectFunction("auto", mhAutoMode)

				c.fixtures["mh1"].MustSelectFunction("color", mhColor)
				c.fixtures["mh2"].MustSelectFunction("color", mhColor)
				c.fixtures["mh3"].MustSelectFunction("color", mhColor)
				c.fixtures["mh4"].MustSelectFunction("color", mhColor)

				c.fixtures["mh1"].MustSelectFunction("gobo", mhGobo)
				c.fixtures["mh2"].MustSelectFunction("gobo", mhGobo)
				c.fixtures["mh3"].MustSelectFunction("gobo", mhGobo)
				c.fixtures["mh4"].MustSelectFunction("gobo", mhGobo)

				c.fixtures["mh1"].MustSetValue("pan", "pan", mhPan)
				c.fixtures["mh2"].MustSetValue("pan", "pan", mhPan)
				c.fixtures["mh3"].MustSetValue("pan", "pan", mhPan)
				c.fixtures["mh4"].MustSetValue("pan", "pan", mhPan)

				c.fixtures["mh1"].MustSetValue("tilt", "tilt", mhTilt)
				c.fixtures["mh2"].MustSetValue("tilt", "tilt", mhTilt)
				c.fixtures["mh3"].MustSetValue("tilt", "tilt", mhTilt)
				c.fixtures["mh4"].MustSetValue("tilt", "tilt", mhTilt)

				//c.fixtures["mh1"].MustSetValue("auto", "sound", mhAutoModeRaw)
				//c.fixtures["mh2"].MustSetValue("auto", "sound", mhAutoModeRaw)
				//c.fixtures["mh3"].MustSetValue("auto", "sound", mhAutoModeRaw)
				//c.fixtures["mh4"].MustSetValue("auto", "sound", mhAutoModeRaw)

				//c.fixtures["sl"].MustSelectFunction("move-macro", macro)
				//c.fixtures["sr"].MustSelectFunction("move-macro", macro)
				//c.fixtures["sl"].MustSetValue("move-macro-speed", "speed", macroSpeed)
				//c.fixtures["sr"].MustSetValue("move-macro-speed", "speed", macroSpeed)

				//c.fixtures["sl"].MustSelectFunction("settings", "blackout-gobo-on")
				//c.fixtures["sr"].MustSelectFunction("settings", "blackout-gobo-on")

				//c.fixtures["sl"].MustSetValue("strobe", "pulse", 231)
				//c.fixtures["sr"].MustSetValue("strobe", "pulse", 231)

				//c.effects["strobe"].Apply()
				//c.effects["rl1"].Apply()

				// Apply active scene
				if c.activeScene != "" {
					c.scenes[c.activeScene].Apply()
				}

				// Apply active effects
				var statusEffects []string
				for name, active := range c.activeEffects {
					if active {
						c.effects[name].Apply()
						statusEffects = append(statusEffects, name)
					}
				}

				// Update status bar
				sort.Strings(statusEffects)
				bar.Describe(fmt.Sprintf("[green]Running [yellow][%s] [cyan]%v [reset]",
					c.activeScene, statusEffects))

				dmx.FrameLock.Unlock()

				execTime := time.Since(start)
				if execTime > c.tickRate {
					fmt.Println("out of time: ", execTime)
				}
			}
		}
	}()
	return nil
}

func (c *Controller) SetActiveScene(in string) error {
	if _, ok := c.scenes[in]; !ok {
		return fmt.Errorf("unknown scene %q", in)
	}
	c.activeScene = in
	return nil
}

func (c *Controller) ToggleEffect(in string) error {
	if _, ok := c.effects[in]; !ok {
		return fmt.Errorf("unknown effect %q", in)
	}
	running := c.activeEffects[in]
	if running {
		c.effects[in].Stop()
		c.effects[in].Apply()
	} else {
		c.effects[in].Start()
	}
	c.activeEffects[in] = !running
	return nil
}

type Universe struct {
	name   string
	output *dmx.Frame
}

func (u *Universe) GetOutput() *dmx.Frame {
	return u.output
}
