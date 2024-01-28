package show

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"go.skymyer.dev/show-control/app/logger"
	"go.skymyer.dev/show-control/common"
	"go.skymyer.dev/show-control/config"
	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/dmx/driver"
	"go.skymyer.dev/show-control/dmx/driver/artnet"
	"go.skymyer.dev/show-control/library"
)

func NewFromConfig(file string, opts ...ControllerOpt) (*Controller, error) {
	var setup = &config.Setup{}
	if err := common.LoadFromFile(file, setup); err != nil {
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
		tickRate:  22 * time.Millisecond,
		setup:     setup,
		universes: make(map[string]*Universe),
		fixtures:  make(map[string]*Fixture),
		devices:   make(map[string]driver.Driver),
		executors: make(map[string]*Executor),
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("option error: %v", err)
		}
	}

	// Initialize Art-Net if enabled
	if setup.Artnet.Enabled {
		var err error
		artnet.DefaultController, err = artnet.NewController(setup.Artnet.Network)
		if err != nil {
			return nil, fmt.Errorf("Art-Net initialization error: %v", err)
		}
		artnet.DefaultController.Run()
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
		device.SetUniverse(universe.Output.Universe, output)

		// Register the universe
		c.universes[name] = &Universe{
			name:   name,
			output: output,
		}
	}

	// Setup fixtures
	for name, conf := range setup.Fixtures {
		// skip virtual devices, will be initialized in second iteration
		if conf.Kind == "virtual" {
			continue
		}
		if _, ok := c.universes[conf.Universe]; !ok {
			return nil, fmt.Errorf("unknown universe assignment %q for fixture %q", conf.Universe, name)
		}
		f, err := NewFixture(c.library, name, conf, c.universes[conf.Universe].GetOutput())
		if err != nil {
			return nil, fmt.Errorf("load fixture %q: %v", name, err)
		}
		c.fixtures[name] = f
	}

	// Setup virtual fixtures
	for name, conf := range setup.Fixtures {
		if conf.Kind != "virtual" {
			continue
		}
		if _, ok := c.fixtures[conf.Real]; !ok {
			return nil, fmt.Errorf("unknown real fixture %q for virtual fixture %q", conf.Real, name)
		}
		f, err := NewVirtualFixture(c.library, name, conf, setup.Fixtures[conf.Real], c.fixtures[conf.Real])
		if err != nil {
			return nil, fmt.Errorf("load virtual fixture %q: %v", name, err)
		}
		c.fixtures[name] = f
	}

	return c, nil
}

type Controller struct {
	shutdownCh chan bool
	setup      *config.Setup
	library    library.Fixtures
	tickRate   time.Duration
	devices    map[string]driver.Driver
	universes  map[string]*Universe
	fixtures   map[string]*Fixture
	executors  map[string]*Executor
	execLock   sync.Mutex
}

func (c *Controller) Stop() error {
	if c.shutdownCh != nil {
		c.shutdownCh <- true
	}
	return nil
}

func (c *Controller) Run() error {

	// Startup device output
	for name, device := range c.devices {
		logger.Default.Debug("starting device", zap.String("device", name))
		device.Run()
	}

	// Main show controller handler
	ticker := time.NewTicker(c.tickRate)
	c.shutdownCh = make(chan bool)

	go func() {
		defer ticker.Stop()
		defer logger.Default.Debug("show controller terminated")

		logger.Default.Debug("start")

		for {
			select {
			case <-c.shutdownCh:
				c.shutdownCh = nil
				logger.Default.Debug("terminating")

				// Terminate executors
				for name := range c.executors {
					c.TerminateExecutor(name)
				}

				// Stop all devices
				for name, device := range c.devices {
					logger.Default.Debug("stopping device", zap.String("device", name))
					device.Stop()
				}

				// Stop ArtNet controller if initialized
				if artnet.DefaultController != nil {
					artnet.DefaultController.Stop()
				}

				logger.Default.Debug("stopped")
				return
			case t := <-ticker.C:
				start := time.Now()

				// Apply executors
				c.execLock.Lock()
				for _, exec := range c.executors {
					for _, ff := range exec.renderFixtureFeatures(t) {
						c.setFixtureFeature(ff)
					}
				}

				// Apply universe DMX frames
				for _, u := range c.universes {
					if len(c.executors) > 0 {
						u.ApplyAndClear()
					} else {
						u.Apply()
					}
				}
				c.execLock.Unlock()

				execTime := time.Since(start)
				if execTime > c.tickRate {
					logger.Default.Warn("show exec out of range",
						zap.Duration("actual", execTime),
						zap.Duration("expected", c.tickRate),
					)
				}
			}
		}
	}()
	return nil
}

func (c *Controller) SetFixtureFeature(f Feature) {
	c.execLock.Lock()
	defer c.execLock.Unlock()
	c.setFixtureFeature(f)
}

func (c *Controller) setFixtureFeature(f Feature) {
	errs := c.fixtures[f.Fixture].SetFeature(f.Name, f.Config)
	for _, err := range errs {
		logger.Default.Error("show control: SetFeature",
			zap.String("fixture", f.Fixture),
			zap.String("feature", f.Name),
			zap.Error(err),
		)
	}
}

func (c *Controller) NewExecutor(name string) (*Executor, error) {
	c.execLock.Lock()
	defer c.execLock.Unlock()
	if _, exists := c.executors[name]; exists {
		return nil, fmt.Errorf("executor %q already present", name)
	}
	e, err := newExecutor(c.fixtures, c.setup.Groups)
	if err != nil {
		return nil, err
	}
	c.executors[name] = e
	logger.Default.Debug("executor added", zap.String("name", name))
	return e, nil
}

func (c *Controller) TerminateExecutor(name string) error {
	c.execLock.Lock()
	defer c.execLock.Unlock()
	if _, exists := c.executors[name]; !exists {
		return fmt.Errorf("unknown executor %q", name)
	}
	delete(c.executors, name)
	logger.Default.Debug("executor removed", zap.String("name", name))
	return nil
}
