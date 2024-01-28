package console

import (
	"fmt"
	ios "io"
	"os"
	"sync"

	"go.skymyer.dev/show-control/io"
)

func init() {
	io.Register(DRIVER_NAME, NewConsole)
}

const (
	DRIVER_NAME = "console"
)

func NewConsole(device string) (io.Driver, error) {
	return &Console{
		out: os.Stdout,
	}, nil
}

type Console struct {
	out       ios.Writer
	input     chan<- io.InputEvent
	eventLock sync.Mutex
}

func (c *Console) Open(input chan<- io.InputEvent) error {
	c.input = input
	return nil
}

func (c *Console) Handle(event ...interface{}) {
	c.eventLock.Lock()
	defer c.eventLock.Unlock()

	fmt.Fprintf(c.out, "Event received %#v\n", event)
}

func (c *Console) Close() error {
	return nil
}

func (c *Console) GetDevices() []string {
	return []string{}
}

func (c *Console) GetVersion() string {
	return "console"
}

func (c *Console) Sleep() error {
	return nil
}
func (c *Console) Wakeup() error {
	return nil
}
