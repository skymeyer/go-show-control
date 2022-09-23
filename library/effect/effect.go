package effect

import (
	"time"

	"go.skymyer.dev/show-control/library/effect/sine"
)

type Modulator interface {
	Init(bpm uint8) error
	Render(t time.Time, shift float64) float64
}

func New(kind string, c RawConfig) Modulator {
	switch kind {
	case sine.NAME:
		effect := &sine.Effect{}
		c.Unmarshal(effect)
		return effect
	default:
		panic("unknown effect")
	}
}

type RawConfig struct {
	unmarshal func(interface{}) error
}

func (c *RawConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	c.unmarshal = unmarshal
	return nil
}

func (c *RawConfig) Unmarshal(v interface{}) error {
	return c.unmarshal(v)
}
