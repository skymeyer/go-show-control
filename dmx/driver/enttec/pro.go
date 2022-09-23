// Enttec Pro DMX Manual:
// https://cdn.enttec.com/pdf/assets/70304/70304_DMX_USB_PRO_USER_MANUAL.pdf
//
// Default settings as reported by EMU (using firmware 1.44)
//
//	Break time 			   = 96 usec
//	Mark After Break (MAB) = 10 usec
//	Packet Refresh Rate    = 40 per second (or 25 ms)
package enttec

import (
	"fmt"
	"io"
	"time"

	"github.com/tarm/serial"
	"go.uber.org/zap"

	"go.skymyer.dev/show-control/app/logger"
	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/dmx/driver"
)

func init() {
	driver.Register(DRIVER_NAME, NewDMXPro)
}

const (
	DRIVER_NAME     = "enttec-dmx-usb-pro"
	BAUD_RATE       = 115200
	START_DELIMITER = 0x7E
	END_DELIMITER   = 0xE7
	REFRESH_RATE_MS = 25
)

const (
	REPROGRAM_FIRMWARE MessageLabel = iota + 1
	PROGRAM_FLASH_PAGE
	GET_WIDGET_PARAMS
	SET_WIDGET_PARAMS
	RECEIVE_DMX_PACKET
	SEND_DMX_PACKET
	SEND_RDM_PACKET
	RECEIVE_DMX_ON_CHANGE
	RECEIVE_DMX_STATE_CHANGE
	GET_WIDGET_SERIAL
	SEND_RDM_DISCOVERY
)

type MessageLabel int

var (
	ErrDeviceNotInitialized = fmt.Errorf("device not initialized")
	ErrUniverseIndex        = fmt.Errorf("only one universe is supported")
)

func NewDMXPro(device string) (driver.Driver, error) {
	d := &DMXPro{
		device: device,
	}
	if err := d.Open(); err != nil {
		return nil, err
	}
	return d, nil
}

type DMXPro struct {
	device     string
	output     *dmx.Frame
	port       io.ReadWriteCloser
	shutdownCh chan bool
}

func (d *DMXPro) SetUniverse(universe int, output *dmx.Frame) error {
	if universe != 0 {
		return ErrUniverseIndex
	}
	d.output = output
	return nil
}

func (d *DMXPro) Open() error {
	port, err := serial.OpenPort(&serial.Config{
		Name: d.device,
		Baud: BAUD_RATE,
	})
	if err != nil {
		return err
	}
	d.port = port
	return nil
}

func (d *DMXPro) Close() error {
	return d.port.Close()
}

func (d *DMXPro) Run() error {
	tickRate := REFRESH_RATE_MS * time.Millisecond
	ticker := time.NewTicker(tickRate)
	d.shutdownCh = make(chan bool)
	go func() {
		defer ticker.Stop()
		defer logger.Default.Debug("enttec dmx pro terminated")
		for {
			select {
			case <-d.shutdownCh:
				d.shutdownCh = nil
				return
			case <-ticker.C:
				start := time.Now()
				d.send(SEND_DMX_PACKET, d.output)
				execTime := time.Since(start)
				if execTime > tickRate {
					logger.Default.Warn("entec apply exec out of range",
						zap.Duration("actual", execTime),
						zap.Duration("expected", tickRate),
					)
				}
			}
		}
	}()
	return nil
}

func (d *DMXPro) Stop() error {
	if d.shutdownCh != nil {
		d.shutdownCh <- true
	}
	if d.port == nil {
		return ErrDeviceNotInitialized
	}
	// clear out all outputs
	d.send(SEND_DMX_PACKET, dmx.NewDMX512Frame())
	return nil
}

func (d *DMXPro) send(label MessageLabel, frame *dmx.Frame) error {
	if d.port == nil {
		return ErrDeviceNotInitialized
	}

	var (
		msg         []byte
		data        = frame.Read()
		dataSizeLSB = byte(len(data) & 0xFF)
		dataSizeMSB = byte(len(data) >> 8 & 0xFF)
	)

	msg = append(msg, START_DELIMITER)
	msg = append(msg, byte(label))
	msg = append(msg, dataSizeLSB)
	msg = append(msg, dataSizeMSB)
	msg = append(msg, data...)
	msg = append(msg, END_DELIMITER)

	_, err := d.port.Write(msg)
	return err
}
