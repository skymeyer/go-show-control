package show

import (
	"go.skymyer.dev/show-control/library/effect"
)

type Effect struct {
	Kind     string           `yaml:"kind,omitempty"`
	Config   effect.RawConfig `yaml:"config,omitempty"`
	Modulate []Modulator      `yaml:"modulate,omitempty"`

	// GrandMA2 effect implementation
	//SingleShot bool
	//Low   int16
	//High  int16
	//Speed int16
	//Phase int8
	//Width int8
}

type Modulator struct {
	Collections []Collection `yaml:"collections,omitempty"`
	Feature     string       `yaml:"feature,omitempty"`
	Properties  []string     `yaml:"properties,omitempty"`
	Min         *uint16      `yaml:"min,omitempty"`
	Max         *uint16      `yaml:"max,omitempty"`
	Shift       float64      `yaml:"shift,omitempty"`
}

type Collection struct {
	Groups   []string `yaml:"groups,omitempty"`
	Fixtures []string `yaml:"fixtures,omitempty"`
}
