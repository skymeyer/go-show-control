package pantilt

import (
	"encoding/binary"
	"fmt"

	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/library/utils"
)

const (
	NAME = "pantilt"
)

type Spec struct {
	Pan      *Axis  `yaml:"pan,omitempty"`
	PanFine  *Axis  `yaml:"panFine,omitempty"`
	Tilt     *Axis  `yaml:"tilt,omitempty"`
	TiltFine *Axis  `yaml:"tiltFine,omitempty"`
	Speed    *Speed `yaml:"speed,omitempty"`
}

type Axis struct {
	Channel   string `yaml:"channel,omitempty"`
	Attribute string `yaml:"attribute,omitempty"`
}

type Speed struct {
	Channel   string `yaml:"channel,omitempty"`
	Attribute string `yaml:"attribute,omitempty"`
}

type Feature struct {
	Pan   *uint16
	Tilt  *uint16
	Speed *uint16
}

func NewHandler(s *Spec) *Handler {
	return &Handler{spec: s}
}

type Handler struct {
	spec *Spec
}

func (h *Handler) Kind() string {
	return NAME
}

func (h *Handler) Render(in interface{}) (ats []dmx.AttributeValue, err error) {
	f, ok := in.(Feature)
	if !ok {
		return nil, fmt.Errorf("%s apply invalid feature config", NAME)
	}

	if h.spec.Pan != nil && f.Pan != nil {
		value := make([]byte, 2)
		binary.LittleEndian.PutUint16(value, *f.Pan)
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Pan.Channel,
			Attribute: h.spec.Pan.Attribute,
			Value:     value[1],
		})
		if h.spec.PanFine != nil {
			ats = append(ats, dmx.AttributeValue{
				Channel:   h.spec.PanFine.Channel,
				Attribute: h.spec.PanFine.Attribute,
				Value:     value[0],
			})
		}
	}

	if h.spec.Tilt != nil && f.Tilt != nil {
		value := make([]byte, 2)
		binary.LittleEndian.PutUint16(value, *f.Tilt)
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Tilt.Channel,
			Attribute: h.spec.Tilt.Attribute,
			Value:     value[1],
		})
		if h.spec.TiltFine != nil {
			ats = append(ats, dmx.AttributeValue{
				Channel:   h.spec.TiltFine.Channel,
				Attribute: h.spec.TiltFine.Attribute,
				Value:     value[0],
			})
		}
	}

	if h.spec.Speed != nil && f.Speed != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Speed.Channel,
			Attribute: h.spec.Speed.Attribute,
			Value:     uint8(*f.Speed / 256),
		})
	}

	return ats, err
}

func (h *Handler) ConfigFromFeatureValues(in utils.FeatureValues) interface{} {
	f := Feature{}
	for k, v := range in {
		switch k {
		case "pan":
			f.Pan = &v
		case "tilt":
			f.Tilt = &v
		case "speed":
			f.Speed = &v
		}
	}
	return f
}
