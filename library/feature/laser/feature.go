package laser

import (
	"fmt"

	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/library/utils"
)

const (
	NAME = "laser"
)

type Spec struct {
	Off     *ChanAttr
	Static  *ChanAttr
	Dynamic *ChanAttr

	Color        *ChanAttr
	Pattern      *ChanAttr
	Size         *ChanAttr
	Angle        *ChanAttr
	Hangle       *ChanAttr
	Vangle       *ChanAttr
	Hpos         *ChanAttr
	Vpos         *ChanAttr
	Line         *ChanAttr
	Scanspeed    *ChanAttr
	Dynamicspeed *ChanAttr
}

type ChanAttr struct {
	Channel   string
	Attribute string
}

type Feature struct {
	Static  bool
	Dynamic bool

	Color        uint8
	Pattern      uint8
	Size         uint8
	Angle        uint8
	HAngle       uint8
	VAngle       uint8
	HPos         uint8
	VPos         uint8
	Line         uint8
	ScanSpeed    uint8
	DynamicSpeed uint8
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

	// Use off channel if exists an neither static or dymamic is selected
	if !f.Static && !f.Dynamic && h.spec.Off != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Off.Channel,
			Attribute: h.spec.Off.Attribute,
			Value:     0,
		})
	}

	// Enable static or dymamic mode
	if f.Static && h.spec.Static != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Static.Channel,
			Attribute: h.spec.Static.Attribute,
			Value:     0,
		})
	} else if f.Dynamic && h.spec.Dynamic != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Dynamic.Channel,
			Attribute: h.spec.Dynamic.Attribute,
			Value:     0,
		})
	}

	if h.spec.Color != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Color.Channel,
			Attribute: h.spec.Color.Attribute,
			Value:     f.Color,
		})
	}
	if h.spec.Pattern != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Pattern.Channel,
			Attribute: h.spec.Pattern.Attribute,
			Value:     f.Pattern,
		})
	}
	if h.spec.Size != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Size.Channel,
			Attribute: h.spec.Size.Attribute,
			Value:     f.Size,
		})
	}
	if h.spec.Angle != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Angle.Channel,
			Attribute: h.spec.Angle.Attribute,
			Value:     f.Angle,
		})
	}
	if h.spec.Hangle != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Hangle.Channel,
			Attribute: h.spec.Hangle.Attribute,
			Value:     f.HAngle,
		})
	}
	if h.spec.Vangle != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Vangle.Channel,
			Attribute: h.spec.Vangle.Attribute,
			Value:     f.VAngle,
		})
	}
	if h.spec.Hpos != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Hpos.Channel,
			Attribute: h.spec.Hpos.Attribute,
			Value:     f.HPos,
		})
	}
	if h.spec.Vpos != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Vpos.Channel,
			Attribute: h.spec.Vpos.Attribute,
			Value:     f.VPos,
		})
	}
	if h.spec.Line != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Line.Channel,
			Attribute: h.spec.Line.Attribute,
			Value:     f.Line,
		})
	}
	if h.spec.Scanspeed != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Scanspeed.Channel,
			Attribute: h.spec.Scanspeed.Attribute,
			Value:     f.ScanSpeed,
		})
	}
	if h.spec.Dynamicspeed != nil {
		ats = append(ats, dmx.AttributeValue{
			Channel:   h.spec.Dynamicspeed.Channel,
			Attribute: h.spec.Dynamicspeed.Attribute,
			Value:     f.DynamicSpeed,
		})
	}

	return ats, err
}

func (h *Handler) ConfigFromFeatureValues(in utils.FeatureValues) interface{} {
	f := Feature{}
	for k, v := range in {
		switch k {
		case "size":
			f.Size = uint8(v)
		case "angle":
			f.Angle = uint8(v)
		case "hangle":
			f.HAngle = uint8(v)
		case "vangle":
			f.VAngle = uint8(v)
		case "hpos":
			f.HPos = uint8(v)
		case "vpos":
			f.VPos = uint8(v)
		case "scanspeed":
			f.ScanSpeed = uint8(v)
		case "dynamicspeed":
			f.DynamicSpeed = uint8(v)
		}
	}
	return f
}
