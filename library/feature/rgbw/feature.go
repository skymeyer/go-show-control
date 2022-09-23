package rgbw

import (
	"fmt"

	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/library/utils"
)

const (
	NAME = "rgbw"
)

type Spec struct {
	Red   []*RGB
	Green []*RGB
	Blue  []*RGB
	White []*RGB
}

type RGB struct {
	Channel   string
	Attribute string
}

type Feature struct {
	Red   *uint8
	Green *uint8
	Blue  *uint8
	White *uint8
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
	if f.Red != nil {
		for _, rgb := range h.spec.Red {
			ats = append(ats, dmx.AttributeValue{
				Channel:   rgb.Channel,
				Attribute: rgb.Attribute,
				Value:     *f.Red,
			})
		}
	}
	if f.Green != nil {
		for _, rgb := range h.spec.Green {
			ats = append(ats, dmx.AttributeValue{
				Channel:   rgb.Channel,
				Attribute: rgb.Attribute,
				Value:     *f.Green,
			})
		}
	}
	if f.Blue != nil {
		for _, rgb := range h.spec.Blue {
			ats = append(ats, dmx.AttributeValue{
				Channel:   rgb.Channel,
				Attribute: rgb.Attribute,
				Value:     *f.Blue,
			})
		}
	}
	if f.White != nil {
		for _, rgb := range h.spec.White {
			ats = append(ats, dmx.AttributeValue{
				Channel:   rgb.Channel,
				Attribute: rgb.Attribute,
				Value:     *f.White,
			})
		}
	}
	return ats, nil
}

func (h *Handler) ConfigFromFeatureValues(in utils.FeatureValues) interface{} {
	f := Feature{}
	for k, v := range in {
		vbyte := uint8(v)
		switch k {
		case "red":
			f.Red = &vbyte
		case "green":
			f.Green = &vbyte
		case "blue":
			f.Blue = &vbyte
		case "white":
			f.White = &vbyte
		}
	}
	return f
}
