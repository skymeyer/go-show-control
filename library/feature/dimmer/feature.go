package dimmer

import (
	"fmt"

	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/library/utils"
)

const (
	NAME = "dimmer"
)

type Spec struct {
	Channel   string
	Attribute string
}

type Feature struct {
	Value *uint16
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

	if f.Value == nil {
		return ats, nil
	}

	ats = append(ats, dmx.AttributeValue{
		Channel:   h.spec.Channel,
		Attribute: h.spec.Attribute,
		Value:     uint8(*f.Value / 256),
	})
	return ats, err
}

func (h *Handler) ConfigFromFeatureValues(in utils.FeatureValues) interface{} {
	f := Feature{}
	for k, v := range in {
		switch k {
		case "value":
			f.Value = &v
		}
	}
	return f
}
