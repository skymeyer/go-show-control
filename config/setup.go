package config

import (
	"gopkg.in/yaml.v3"
)

type Setup struct {
	Name      string
	Devices   map[string]*Device
	Universes map[string]*Universe
	Fixtures  map[string]*Fixture
	Groups    map[string]*Group
	Scenes    map[string][]*Scene
	Effects   map[string][]*Effect
	HotKeys   map[string]*Hotkey
}

type Device struct {
	Driver string
	Device string
}

type Universe struct {
	Output Output
}

type Output struct {
	Device   string
	Universe int
}

type Fixture struct {
	Kind     string
	Mode     string
	Universe string
	Address  int
}

type Group struct {
	Members []string
}

type Scene struct {
	Kind     string
	Fixtures []string
	Groups   []string
	Spec     interface{} `yaml:"-"`
}

const (
	SCENE_KIND_RGB   = "rgb"
	SCENE_KIND_COLOR = "color"
	SCENE_KIND_GOBO  = "gobo"
)

type RGBSpec struct {
	Color string `yaml:"color"`
}

type ColorSpec struct {
	Color string `yaml:"color"`
}

type GoboSpec struct {
	Gobo string `yaml:"gobo"`
}

// Dynamic Unmarshal scene spec
func (s *Scene) UnmarshalYAML(n *yaml.Node) error {
	type S Scene
	type T struct {
		*S   `yaml:",inline"`
		Spec yaml.Node `yaml:"spec"`
	}

	obj := &T{S: (*S)(s)}
	if err := n.Decode(obj); err != nil {
		return err
	}

	switch s.Kind {
	case SCENE_KIND_RGB:
		s.Spec = new(RGBSpec)
	case SCENE_KIND_COLOR:
		s.Spec = new(ColorSpec)
	case SCENE_KIND_GOBO:
		s.Spec = new(GoboSpec)
	default:
		panic("unknown scene kind")
	}
	return obj.Spec.Decode(s.Spec)
}

type Effect struct {
	Kind     string
	Fixtures []string
	Groups   []string
	Spec     interface{} `yaml:"-"`
}

const (
	EFFECT_KIND_STROBE    = "strobe"
	EFFECT_KIND_MATRIX    = "matrix"
	EFFECT_KIND_PAN_TILT  = "pan_tilt"
	EFFECT_KIND_RAW_FUNC  = "function"
	EFFECT_KIND_RAW_VALUE = "value"
)

type StrobeSpec struct {
	Speed int `yaml:"speed"`
	Off   int `yaml:"off"`
}

type MatrixSpec struct {
	Matrix  [][]string `yaml:"matrix"`
	Cycle   int        `yaml:"cycle"`
	Step    int        `yaml:"step"`
	Pattern uint8      `yaml:"pattern"`
}

type PantiltSpec struct {
	Fixed *struct {
		Pan  int `yaml:"pan"`
		Tilt int `yaml:"tilt"`
	}
}

type RawFunctionEffectSpec struct {
	Channel       string `yaml:"channel"`
	FunctionStart string `yaml:"start"`
	FunctionStop  string `yaml:"stop"`
}

type RawValueEffectSpec struct {
	Channel  string `yaml:"channel"`
	Function string `yaml:"function"`
	Value    int    `yaml:"value"`
}

// Dynamic Unmarshal effect spec
func (e *Effect) UnmarshalYAML(n *yaml.Node) error {
	type E Effect
	type T struct {
		*E   `yaml:",inline"`
		Spec yaml.Node `yaml:"spec"`
	}

	obj := &T{E: (*E)(e)}
	if err := n.Decode(obj); err != nil {
		return err
	}

	switch e.Kind {
	case EFFECT_KIND_STROBE:
		e.Spec = new(StrobeSpec)
	case EFFECT_KIND_MATRIX:
		e.Spec = new(MatrixSpec)
	case EFFECT_KIND_PAN_TILT:
		e.Spec = new(PantiltSpec)
	case EFFECT_KIND_RAW_FUNC:
		e.Spec = new(RawFunctionEffectSpec)
	case EFFECT_KIND_RAW_VALUE:
		e.Spec = new(RawValueEffectSpec)
	default:
		panic("unknown effect kind")
	}
	return obj.Spec.Decode(e.Spec)
}

const (
	HOTKEY_KIND_SCENE  HotkeyKind = "scene"
	HOTKEY_KIND_EFFECT HotkeyKind = "effect"
)

type HotkeyKind string

type Hotkey struct {
	Kind  HotkeyKind
	Value string
}
