package show

import (
	"go.skymyer.dev/show-control/library/feature"
)

type Feature struct {
	Name    string
	Kind    string
	Fixture string
	Config  interface{}
}

type Cue struct {
	Feature  string            `yaml:"feature,omitempty"`
	Groups   []string          `yaml:"groups,omitempty"`
	Fixtures []string          `yaml:"fixtures,omitempty"`
	Config   feature.RawConfig `yaml:"config,omitempty"`
}

type Show struct {
	Cues      map[string][]Cue     `yaml:"cues,omitempty"`
	Effects   map[string][]Effect  `yaml:"effects,omitempty"`
	Sequences map[string]*Sequence `yaml:"sequences,omitempty"`
}
