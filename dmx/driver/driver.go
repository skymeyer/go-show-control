package driver

import (
	"fmt"

	"go.skymyer.dev/show-control/dmx"
)

var drivers = map[string]Factory{}

type Factory func(device string) (Driver, error)

func Register(name string, f Factory) {
	drivers[name] = f
}

func New(driver, device string) (Driver, error) {
	if factory, found := drivers[driver]; found {
		return factory(device)
	}
	return nil, fmt.Errorf("unknown driver %s", driver)
}

type Driver interface {
	Open() error
	Close() error
	Run() error
	Stop() error
	SetUniverse(universe int, output *dmx.Frame) error
}

//type Config struct {
//	Model          string
//	Serial         string
//	Firmware       string
//	BreakTime      int
//	MarkAfterBreak int
//	Rate           int
//}
