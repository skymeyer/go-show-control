package library

import (
	"gopkg.in/yaml.v3"

	"go.skymyer.dev/show-control/common"
	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/library/feature"
)

func FixtureLibraryFromConfig(file string) (Fixtures, error) {
	var fixtures = make(map[string]*Fixture)
	if err := common.LoadFromFile(file, &fixtures); err != nil {
		return nil, err
	}

	// TODO: add validators:: value/ranges, head refs, features, ...

	return fixtures, nil
}

type Fixtures map[string]*Fixture

type Fixture struct {
	Manufacturer string
	Model        string
	Modes        map[string]*Mode
}

type Mode struct {
	Dmx      dmx.Channels
	Features map[string]*Feature
}

type Feature struct {
	Kind string
	Spec interface{} `yaml:"-"`
}

// Dynamic unmarshal feature spec
func (f *Feature) UnmarshalYAML(n *yaml.Node) error {
	type F Feature
	type T struct {
		*F   `yaml:",inline"`
		Spec yaml.Node `yaml:"spec"`
	}

	obj := &T{F: (*F)(f)}
	if err := n.Decode(obj); err != nil {
		return err
	}

	f.Spec = feature.NewSpec(f.Kind)
	return obj.Spec.Decode(f.Spec)
}

type RawConfig struct {
	unmarshal func(interface{}) error
}

func (msg *RawConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	msg.unmarshal = unmarshal
	return nil
}

func (msg *RawConfig) Unmarshal(v interface{}) error {
	return msg.unmarshal(v)
}
