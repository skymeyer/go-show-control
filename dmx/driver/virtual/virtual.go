package virtual

import (
	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/dmx/driver"
)

func init() {
	driver.Register(DRIVER_NAME, NewVirtual)
}

const (
	DRIVER_NAME = "virtual"
)

func NewVirtual(device string) (driver.Driver, error) {
	d := &Virtual{
		device: device,
		output: make(map[int]*dmx.Frame),
	}
	if err := d.Open(); err != nil {
		return nil, err
	}
	return d, nil
}

type Virtual struct {
	device string
	output map[int]*dmx.Frame
}

func (d *Virtual) SetUniverse(universe int, output *dmx.Frame) error {
	d.output[universe] = output
	return nil
}

func (d *Virtual) Open() error {
	return nil
}

func (d *Virtual) Close() error {
	return nil
}

func (d *Virtual) Run() error {
	return nil
}

func (d *Virtual) Stop() error {
	return nil
}
