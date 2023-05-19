package feature

import (
	"go.skymyer.dev/show-control/dmx"

	"go.skymyer.dev/show-control/library/feature/cannon"
	"go.skymyer.dev/show-control/library/feature/dimmer"
	"go.skymyer.dev/show-control/library/feature/gobo"
	"go.skymyer.dev/show-control/library/feature/laser"
	"go.skymyer.dev/show-control/library/feature/pantilt"
	"go.skymyer.dev/show-control/library/feature/rgbw"
	"go.skymyer.dev/show-control/library/feature/strobe"
	"go.skymyer.dev/show-control/library/utils"
)

type Handler interface {
	// Render take a feature configuration and outputs the DMX attribute values
	Render(interface{}) ([]dmx.AttributeValue, error)
	ConfigFromFeatureValues(utils.FeatureValues) interface{}
	Kind() string
}

func NewHandler(feature string, spec interface{}, channels dmx.Channels) Handler {
	switch feature {
	case cannon.NAME:
		s := spec.(*cannon.Spec)
		return cannon.NewHandler(s, channels)
	case dimmer.NAME:
		s := spec.(*dimmer.Spec)
		return dimmer.NewHandler(s)
	case gobo.NAME:
		s := spec.(*gobo.Spec)
		return gobo.NewHandler(s, channels)
	case laser.NAME:
		s := spec.(*laser.Spec)
		return laser.NewHandler(s)
	case pantilt.NAME:
		s := spec.(*pantilt.Spec)
		return pantilt.NewHandler(s)
	case rgbw.NAME:
		s := spec.(*rgbw.Spec)
		return rgbw.NewHandler(s)
	case strobe.NAME:
		s := spec.(*strobe.Spec)
		return strobe.NewHandler(s, channels)
	default:
		panic("unknown feature handler")
	}
}

func NewSpec(kind string) interface{} {
	switch kind {
	case cannon.NAME:
		return new(cannon.Spec)
	case dimmer.NAME:
		return new(dimmer.Spec)
	case gobo.NAME:
		return new(gobo.Spec)
	case laser.NAME:
		return new(laser.Spec)
	case pantilt.NAME:
		return new(pantilt.Spec)
	case rgbw.NAME:
		return new(rgbw.Spec)
	case strobe.NAME:
		return new(strobe.Spec)
	default:
		panic("unknown feature spec")
	}
}

func NewConfig(kind string, c RawConfig) interface{} {
	switch kind {
	case cannon.NAME:
		config := cannon.Feature{}
		c.Unmarshal(&config)
		return config
	case dimmer.NAME:
		config := dimmer.Feature{}
		c.Unmarshal(&config)
		return config
	case gobo.NAME:
		config := gobo.Feature{}
		c.Unmarshal(&config)
		return config
	case laser.NAME:
		config := laser.Feature{}
		c.Unmarshal(&config)
		return config
	case pantilt.NAME:
		config := pantilt.Feature{}
		c.Unmarshal(&config)
		return config
	case rgbw.NAME:
		config := rgbw.Feature{}
		c.Unmarshal(&config)
		return config
	case strobe.NAME:
		config := strobe.Feature{}
		c.Unmarshal(&config)
		return config
	default:
		panic("unknown feature config")
	}
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
