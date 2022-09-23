package strobe

import (
	"fmt"

	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/library/utils"
)

const (
	NAME = "strobe"
)

type Spec struct {
	Channel   string
	Attribute string
	Off       *OffSpec
}

type OffSpec struct {
	Channel   string
	Attribute string
	Value     uint8
}

type Feature struct {
	Value uint8
}

func NewHandler(s *Spec, channels dmx.Channels) *Handler {
	return &Handler{
		spec: s,
	}
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

	if h.spec.Off != nil && f.Value == 0 {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Off.Channel,
			Attribute: h.spec.Off.Attribute,
			Value:     0,
		})
	} else {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Channel,
			Attribute: h.spec.Attribute,
			Value:     f.Value,
		})
	}
	return ats, err
}

func (h *Handler) ConfigFromFeatureValues(in utils.FeatureValues) interface{} {
	f := Feature{}
	for k, v := range in {
		switch k {
		case "value":
			f.Value = uint8(v)
		}
	}
	return f
}
