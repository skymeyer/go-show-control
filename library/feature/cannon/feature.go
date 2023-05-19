package cannon

import (
	"fmt"

	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/library/utils"
)

const (
	NAME = "cannon"
)

type Spec struct {
	Cannon1 *Cannon
	Cannon2 *Cannon
	Cannon3 *Cannon
	Cannon4 *Cannon
}

type Cannon struct {
	Channel string
	Standby string
	Launch  string
}

type Feature struct {
	Launch1 bool
	Launch2 bool
	Launch3 bool
	Launch4 bool
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

	if h.spec.Cannon1 != nil {
		if ch, ok := h.channels[h.spec.Cannon1.Channel]; ok {
			attr := h.spec.Cannon1.Standby
			if f.Launch1 {
				attr = h.spec.Cannon1.Launch
			}
			if value, ok := ch.Attributes[attr]; ok {
				ats = append(ats, dmx.AttributeValue{
					Channel:   h.spec.Cannon1.Channel,
					Attribute: attr,
					Value:     value.Min,
				})
			}
		}
	}

	if h.spec.Cannon2 != nil {
		if ch, ok := h.channels[h.spec.Cannon2.Channel]; ok {
			attr := h.spec.Cannon2.Standby
			if f.Launch2 {
				attr = h.spec.Cannon2.Launch
			}
			if value, ok := ch.Attributes[attr]; ok {
				ats = append(ats, dmx.AttributeValue{
					Channel:   h.spec.Cannon2.Channel,
					Attribute: attr,
					Value:     value.Min,
				})
			}
		}
	}

	if h.spec.Cannon3 != nil {
		if ch, ok := h.channels[h.spec.Cannon3.Channel]; ok {
			attr := h.spec.Cannon3.Standby
			if f.Launch3 {
				attr = h.spec.Cannon3.Launch
			}
			if value, ok := ch.Attributes[attr]; ok {
				ats = append(ats, dmx.AttributeValue{
					Channel:   h.spec.Cannon3.Channel,
					Attribute: attr,
					Value:     value.Min,
				})
			}
		}
	}

	if h.spec.Cannon4 != nil {
		if ch, ok := h.channels[h.spec.Cannon4.Channel]; ok {
			attr := h.spec.Cannon4.Standby
			if f.Launch4 {
				attr = h.spec.Cannon4.Launch
			}
			if value, ok := ch.Attributes[attr]; ok {
				ats = append(ats, dmx.AttributeValue{
					Channel:   h.spec.Cannon4.Channel,
					Attribute: attr,
					Value:     value.Min,
				})
			}
		}
	}

	return ats, err
}

func (h *Handler) ConfigFromFeatureValues(in utils.FeatureValues) interface{} {
	return Feature{}
}
