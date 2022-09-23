package gobo

import (
	"fmt"

	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/library/utils"
)

const (
	NAME = "gobo"
)

type Spec struct {
	Channel string
}

type Feature struct {
	Gobo string
}

func NewHandler(s *Spec, channels dmx.Channels) *Handler {
	return &Handler{
		spec:     s,
		channels: channels,
	}
}

type Handler struct {
	spec     *Spec
	channels dmx.Channels
}

func (h *Handler) Kind() string {
	return NAME
}

func (h *Handler) Render(in interface{}) (ats []dmx.AttributeValue, err error) {
	f, ok := in.(Feature)
	if !ok {
		return nil, fmt.Errorf("%s apply invalid feature config", NAME)
	}

	// Lookup attribute value
	if ch, ok := h.channels[h.spec.Channel]; ok {
		if value, ok := ch.Attributes[f.Gobo]; ok {
			ats = append(ats, dmx.AttributeValue{
				Channel:   h.spec.Channel,
				Attribute: f.Gobo,
				Value:     value.Min,
			})
			return ats, nil
		}
	}

	return ats, fmt.Errorf("cannot find gobo %q on channel %q", f.Gobo, h.spec.Channel)
}

func (h *Handler) ConfigFromFeatureValues(in utils.FeatureValues) interface{} {
	return Feature{}
}
