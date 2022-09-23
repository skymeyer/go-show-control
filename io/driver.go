package io

import (
	"fmt"
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
	Open(chan<- InputEvent) error
	Close() error
	Handle(event ...interface{})

	GetDevices() []string
	GetVersion() string
	Sleep() error
	Wakeup() error
}
