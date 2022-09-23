package show

import (
	"go.skymyer.dev/show-control/dmx"
)

type Universe struct {
	name   string
	output *dmx.Frame
}

func (u *Universe) Apply() {
	u.output.Apply(false)
}

func (u *Universe) ApplyAndClear() {
	u.output.Apply(true)
}

func (u *Universe) GetOutput() *dmx.Frame {
	return u.output
}
