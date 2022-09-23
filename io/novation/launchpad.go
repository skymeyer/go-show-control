package novation

import (
	"bytes"
	"fmt"
	"math"
	"sync"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"

	"go.skymyer.dev/show-control/io"
)

func init() {
	io.Register(DRIVER_NAME, NewLaunchpadMiniMk3)
}

const (
	DRIVER_NAME        = "novation-launchpad-mini-mk3"
	LAUNCHPAD_MINI_MK3 = "Launchpad Mini MK3 LPMiniMK3 MIDI"
)

var (
	// Universal SysEx Device Inquiry, not using Novation specific header
	sysExDeviceInquiry    = []byte{0x7E, 0x7F, 0x06, 0x01}
	sysExDeviceInquiryApp = []byte{0x7E, 0x00, 0x06, 0x02, 0x00, 0x20, 0x29, 0x13, 0x01, 0x00, 0x00}

	// Novation SysEx message using its own header
	sysExHeader = []byte{0x00, 0x20, 0x29, 0x02, 0x0D}

	// Layout management
	sysExLayout           = []byte{0x00}
	sysExLayoutProgrammer = append(sysExLayout, 0x7F)

	// Mode management
	sysExMode           = []byte{0x0E}
	sysExModeProgrammer = append(sysExMode, 0x01)

	// LED lighting
	sysExLight           = []byte{0x03}
	sysExColorSpecStatic = byte(0x00)
	sysExColorSpecFlash  = byte(0x01)
	sysExColorSpecPulse  = byte(0x02)
	sysExColorSpecRGB    = byte(0x03)

	// Text scrolling
	sysExText       = []byte{0x07}
	sysExTextNoLoop = append(sysExText, 0x00)
	sysExTextLoop   = append(sysExText, 0x01)

	// Power management
	sysExPower      = []byte{0x9}
	sysExPowerSleep = append(sysExPower, 0x00)
	sysExPowerWakup = append(sysExPower, 0x01)

	colorMode = map[io.Mode]byte{
		io.MODE_MAIN: 0x00,
		io.MODE_1:    0x57,
		io.MODE_2:    0x48,
		io.MODE_3:    0x77,
		io.MODE_4:    0x45,
	}

	colorModePage = map[io.Control]byte{
		io.CTR_PAGE_1: colorMode[io.MODE_1],
		io.CTR_PAGE_2: colorMode[io.MODE_2],
		io.CTR_PAGE_3: colorMode[io.MODE_3],
		io.CTR_PAGE_4: colorMode[io.MODE_4],
	}

	colorControl = map[io.Control]byte{
		io.CTR_PAGE_1: 0x66,
		io.CTR_PAGE_2: 0x66,
		io.CTR_PAGE_3: 0x66,
		io.CTR_PAGE_4: 0x66,
	}

	colorControlSelect = map[io.Control]byte{
		io.CTR_PAGE_1: 0x67,
		io.CTR_PAGE_2: 0x67,
		io.CTR_PAGE_3: 0x67,
		io.CTR_PAGE_4: 0x67,
	}

	colorKind = map[io.ButtonKind]byte{
		io.BTN_KIND_OFF:  0x00,
		io.BTN_KIND_STOP: 0x48,

		io.BTN_KIND_DEFAULT:        0x2D,
		io.BTN_KIND_DEFAULT_SELECT: 0x15,

		io.BTN_KIND_COLOR_BLUE:    0x2D,
		io.BTN_KIND_COLOR_CYAN:    0x4E,
		io.BTN_KIND_COLOR_GRAY:    0x01,
		io.BTN_KIND_COLOR_GREEN:   0x57,
		io.BTN_KIND_COLOR_MAGENTA: 0x51,
		io.BTN_KIND_COLOR_ORANGE:  0x6C,
		io.BTN_KIND_COLOR_RED:     0x48,
		io.BTN_KIND_COLOR_WHITE:   0x03,
		io.BTN_KIND_COLOR_YELLOW:  0x4A,

		io.BTN_KIND_CALL:         0x0B,
		io.BTN_KIND_CALL_SELECT:  0x54,
		io.BTN_KIND_ARROW:        0x0B,
		io.BTN_KIND_ARROW_SELECT: 0x54,
	}

	colorKindDual = map[io.ButtonKind][2]byte{
		io.BTN_KIND_DUAL_MAGENTA_CYAN:  {colorKind[io.BTN_KIND_COLOR_MAGENTA], colorKind[io.BTN_KIND_COLOR_CYAN]},
		io.BTN_KIND_DUAL_MAGENTA_BLUE:  {colorKind[io.BTN_KIND_COLOR_MAGENTA], colorKind[io.BTN_KIND_COLOR_BLUE]},
		io.BTN_KIND_DUAL_GREEN_BLUE:    {colorKind[io.BTN_KIND_COLOR_GREEN], colorKind[io.BTN_KIND_COLOR_BLUE]},
		io.BTN_KIND_DUAL_GREEN_YELLOW:  {colorKind[io.BTN_KIND_COLOR_GREEN], colorKind[io.BTN_KIND_COLOR_YELLOW]},
		io.BTN_KIND_DUAL_ORANGE_YELLOW: {colorKind[io.BTN_KIND_COLOR_ORANGE], colorKind[io.BTN_KIND_COLOR_YELLOW]},
		io.BTN_KIND_DUAL_ORANGE_RED:    {colorKind[io.BTN_KIND_COLOR_ORANGE], colorKind[io.BTN_KIND_COLOR_RED]},
	}
)

func init() {
	controlMapReverse = make(map[io.Control]byte)
	for k, v := range controlMap {
		controlMapReverse[v] = k
	}
	buttonMapReverse = make(map[io.Button]byte)
	for k, v := range buttonMap {
		buttonMapReverse[v] = k
	}
}

func NewLaunchpadMiniMk3(device string) (io.Driver, error) {
	return &LaunchpadMiniMk3{
		device:  device,
		welcome: "", // Welcome string
	}, nil
}

type LaunchpadMiniMk3 struct {
	welcome   string
	device    string
	version   string // application version
	mode      byte   // the selected mode
	layout    byte   // the selected layout
	stop      func()
	in        drivers.In
	out       drivers.Out
	send      func(midi.Message) error
	eventLock sync.Mutex
}

func (d *LaunchpadMiniMk3) Open(input chan<- io.InputEvent) error {
	var err error

	// input channel device
	d.in, err = midi.FindInPort(d.device)
	if err != nil {
		return fmt.Errorf("%q find in port : %v", DRIVER_NAME, err)
	}

	// output channel device
	d.out, err = midi.FindOutPort(d.device)
	if err != nil {
		return fmt.Errorf("%q find out port: %v", DRIVER_NAME, err)
	}

	// send wrapper for output channel
	d.send, err = midi.SendTo(d.out)
	if err != nil {
		return fmt.Errorf("send prep: %v", err)
	}

	// start listening on input channel
	d.stop, err = d.listen(input)
	if err != nil {
		return fmt.Errorf("%q listen: %v", DRIVER_NAME, err)
	}

	// poll version
	if err := d.sendRawSysEx(sysExDeviceInquiry); err != nil {
		return fmt.Errorf("device inquiry: %v", err)
	}

	if err := d.sendSysEx(sysExLayoutProgrammer); err != nil {
		return fmt.Errorf("set layout: %v", err)
	}
	if err := d.sendSysEx(sysExModeProgrammer); err != nil {
		return fmt.Errorf("set mode: %v", err)
	}

	// allow some time for sysex exchange
	time.Sleep(1 * time.Second)

	// send welcome message
	if len(d.welcome) > 0 {
		d.sendSysEx(outputText(d.welcome, false))
	}

	return nil
}

func (d *LaunchpadMiniMk3) Close() error {
	if d.stop != nil {
		d.stop()
	}
	if d.in != nil {
		if err := d.in.Close(); err != nil {
			return err
		}
	}
	if d.out != nil {
		if err := d.out.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (d *LaunchpadMiniMk3) Handle(events ...interface{}) {
	d.eventLock.Lock()
	defer d.eventLock.Unlock()

	var gridControl = make(gridControl)
	var gridButton = make(gridButton)

	for _, event := range events {
		switch e := event.(type) {
		case io.ControlChangeEvent:
			// If no partial event, we explicitly set all buttons to off
			if !e.Partial {
				gridControl = emptyControlGrid()
			}

			// Always set logo to selected mode
			gridControl[io.CTR_SELECT] = staticColor(colorMode[e.Mode])

			// Prepare colors for requested control buttons
			for _, btn := range e.Grid {

				// Default color (unselected)
				var color = staticColor(colorControl[btn.Control])

				// Highlight/select page button
				if btn.Control == io.PageToControl[e.Page] {
					color = staticColor(colorControlSelect[btn.Control])
				}

				// When in main mode, we use the main colors for page buttons
				if e.Mode == io.MODE_MAIN {
					if override, ok := colorModePage[btn.Control]; ok {
						color = staticColor(byte(override))
					}
				}

				// Set stop button if available
				if btn.Control == io.CTR_BACK {
					color = staticColor(colorKind[io.BTN_KIND_STOP])
				}

				// add the button color
				gridControl[btn.Control] = color
			}
		case io.ButtonChangeEvent:
			// If no partial event, we explicitly set all buttons to off
			if !e.Partial {
				gridButton = emptyButtonGrid()
			}

			// Set buttons
			for _, btn := range e.Grid {
				switch btn.Kind {
				case io.BTN_KIND_WIDGET_SLIDER_RED:
					for k, v := range createWidgetSlider(btn.Button, RGB_RED) {
						gridButton[k] = v
					}
				case io.BTN_KIND_WIDGET_SLIDER_GREEN:
					for k, v := range createWidgetSlider(btn.Button, RGB_GREEN) {
						gridButton[k] = v
					}
				case io.BTN_KIND_WIDGET_SLIDER_BLUE:
					for k, v := range createWidgetSlider(btn.Button, RGB_BLUE) {
						gridButton[k] = v
					}
				case io.BTN_KIND_WIDGET_SLIDER_WHITE:
					for k, v := range createWidgetSlider(btn.Button, RGB_WHITE) {
						gridButton[k] = v
					}
				default:
					if dual, ok := colorKindDual[btn.Kind]; ok {
						gridButton[btn.Button] = dualColor(dual[0], dual[1])
					} else {
						if color, ok := colorKind[btn.Kind]; ok {
							gridButton[btn.Button] = staticColor(color)
						} else {
							gridButton[btn.Button] = staticColor(colorKind[io.BTN_KIND_DEFAULT])
						}
					}
				}
			}
		}
	}

	// send all buttons to be changed
	lights := append(sysExLight, controlGridtoColorSpec(gridControl)...)
	lights = append(lights, buttonGridtoColorSpec(gridButton)...)
	d.sendSysEx(lights)
}

type RgbColor []byte

var (
	RGB_RED   RgbColor = []byte{0xFF, 0x00, 0x00}
	RGB_GREEN RgbColor = []byte{0x00, 0xFF, 0x00}
	RGB_BLUE  RgbColor = []byte{0x00, 0x00, 0xFF}
	RGB_WHITE RgbColor = []byte{0xFF, 0xFF, 0xFF}
)

func createWidgetSlider(start io.Button, color RgbColor) gridButton {
	var widget = make(gridButton)

	// Find row to draw the widget
	var buttons []io.Button

	for _, list := range buttonGrid {
		if list[0] == start {
			buttons = list
			break
		}
	}

	// If no row is found, silently return
	if buttons == nil {
		return widget
	}

	var weight byte
	for row, button := range buttons {
		weight = byte(math.Pow(2, float64(row)) + float64(weight))
		widget[button] = rgbColor(weight&color[0], weight&color[1], weight&color[2])
	}

	return widget
}

func (d *LaunchpadMiniMk3) GetDevices() (list []string) {
	if d.in != nil {
		list = append(list, d.in.String())
	}
	if d.out != nil {
		list = append(list, d.out.String())
	}
	return list
}

func (d *LaunchpadMiniMk3) GetVersion() string {
	return d.version
}

func (d *LaunchpadMiniMk3) Sleep() error {
	return d.sendSysEx(sysExPowerSleep)
}

func (d *LaunchpadMiniMk3) Wakeup() error {
	return d.sendSysEx(sysExPowerWakup)
}

func (d *LaunchpadMiniMk3) listen(input chan<- io.InputEvent) (func(), error) {
	return midi.ListenTo(d.in, func(msg midi.Message, timestampms int32) {
		var bt []byte
		var ch, key, val uint8
		switch {
		case msg.GetSysEx(&bt):
			d.handleSysEx(bt)
		case msg.GetControlChange(&ch, &key, &val):
			var (
				control io.Control
				button  io.Button
				action  = io.ACTION_RELEASE
			)
			// CC messages can be both control or button actions
			if btn, ok := controlMap[key]; ok {
				control = btn
			}
			if btn, ok := buttonMap[key]; ok {
				button = btn
			}
			if control > 0 || button > 0 {
				if val > 0 {
					action = io.ACTION_PRESS
				}
				input <- io.InputEvent{
					Action:  action,
					Button:  button,
					Control: control,
				}
			}
		case msg.GetNoteStart(&ch, &key, &val):
			if btn, ok := buttonMap[key]; ok {
				input <- io.InputEvent{
					Action: io.ACTION_PRESS,
					Button: btn,
				}
			}
		case msg.GetNoteEnd(&ch, &key):
			if btn, ok := buttonMap[key]; ok {
				input <- io.InputEvent{
					Action: io.ACTION_RELEASE,
					Button: btn,
				}
			}
		default:
			fmt.Printf("unknown message: % X\n", msg.Bytes())
		}
	}, midi.UseSysEx(), midi.SysExBufferSize(4096))
}

func (d *LaunchpadMiniMk3) handleSysEx(msg []byte) {

	// Strip novation specific header if present
	msg = bytes.TrimPrefix(msg, sysExHeader)

	switch {
	// Track version when polled
	case bytes.HasPrefix(msg, sysExDeviceInquiryApp):
		d.version = fmt.Sprintf("%d%d%d%d", msg[11], msg[12], msg[13], msg[14])

	// Track mode change
	case bytes.HasPrefix(msg, sysExMode):
		d.mode = bytes.TrimPrefix(msg, sysExMode)[0]

	// Track layout change
	case bytes.HasPrefix(msg, sysExLayout):
		d.layout = bytes.TrimPrefix(msg, sysExLayout)[0]

	// Log anything unknown (should be an option later on)
	default:
		fmt.Printf("unhandled SysEx message: % X\n", msg)
	}
}

// sendSysEx sends a SysEx message with Novation header.
func (d *LaunchpadMiniMk3) sendSysEx(msg []byte) error {
	return d.sendRawSysEx(append(sysExHeader, msg...))
}

// sendRawSysEx sends a raw SysEx message without Novation header.
func (d *LaunchpadMiniMk3) sendRawSysEx(msg []byte) error {
	return d.send(midi.SysEx(msg))
}

func outputText(txt string, loop bool) []byte {
	var msg []byte
	if loop {
		msg = sysExTextLoop
	} else {
		msg = sysExTextNoLoop
	}

	msg = append(msg, 0x07, 0x00, 0x37)
	msg = append(msg, []byte(txt)...)

	return msg
}

type gridControl map[io.Control][]byte

func emptyControlGrid() gridControl {
	grid := make(gridControl)
	grid[io.CTR_SELECT] = buttonOff()
	grid[io.CTR_PAGE_1] = buttonOff()
	grid[io.CTR_PAGE_2] = buttonOff()
	grid[io.CTR_PAGE_3] = buttonOff()
	grid[io.CTR_PAGE_4] = buttonOff()
	grid[io.CTR_BACK] = buttonOff()
	return grid
}

type gridButton map[io.Button][]byte

func emptyButtonGrid() gridButton {
	grid := make(gridButton)
	grid[io.BTN_ARROW_UP] = buttonOff()
	grid[io.BTN_ARROW_DOWN] = buttonOff()
	grid[io.BTN_ARROW_LEFT] = buttonOff()
	grid[io.BTN_ARROW_RIGHT] = buttonOff()
	grid[io.BTN_CALL_1] = buttonOff()
	grid[io.BTN_CALL_2] = buttonOff()
	grid[io.BTN_CALL_3] = buttonOff()
	grid[io.BTN_CALL_4] = buttonOff()
	grid[io.BTN_CALL_5] = buttonOff()
	grid[io.BTN_CALL_6] = buttonOff()
	grid[io.BTN_CALL_7] = buttonOff()

	// Generate other buttons from grid
	for _, row := range buttonGrid {
		for _, button := range row {
			grid[button] = buttonOff()
		}
	}
	return grid
}

func staticColor(color byte) []byte {
	return []byte{sysExColorSpecStatic, color}
}

func pulseColor(color byte) []byte {
	return []byte{sysExColorSpecPulse, color}
}

func dualColor(color1, color2 byte) []byte {
	return []byte{sysExColorSpecFlash, color1, color2}
}

func rgbColor(r, g, b byte) []byte {
	// Only RGB values 0-127 are allowed so dropping LSB for each
	return []byte{sysExColorSpecRGB, r >> 1, g >> 1, b >> 1}
}

func buttonOff() []byte {
	return staticColor(colorKind[io.BTN_KIND_OFF])
}

func controlGridtoColorSpec(c gridControl) (out []byte) {
	for btn, spec := range c {
		out = append(out, spec[0:1]...)
		out = append(out, controlMapReverse[btn])
		out = append(out, spec[1:]...)
	}
	return out
}

func buttonGridtoColorSpec(c gridButton) (out []byte) {
	for btn, spec := range c {
		out = append(out, spec[0:1]...)
		out = append(out, buttonMapReverse[btn])
		out = append(out, spec[1:]...)
	}
	return out
}
