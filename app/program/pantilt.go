package program

import (
	"go.skymyer.dev/show-control/io"
)

type PanTilt struct {
	ButtonsPan   []io.Button
	ButtonsTilt  []io.Button
	ButtonsSpeed []io.Button
	PanValue     uint16
	TiltValue    uint16
	SpeedValue   uint8
}

func (p *PanTilt) Handle(button io.Button) bool {
	for index, btn := range p.ButtonsPan {
		if btn == button {
			p.PanValue = buttonWeightedUint16(p.ButtonsPan, index)
			return true
		}
	}
	for index, btn := range p.ButtonsTilt {
		if btn == button {
			p.TiltValue = buttonWeightedUint16(p.ButtonsTilt, index)
			return true
		}
	}
	for index, btn := range p.ButtonsSpeed {
		if btn == button {
			p.SpeedValue = buttonWeightedUint8(p.ButtonsSpeed, index)
			return true
		}
	}
	return false
}

func (p *PanTilt) HandleArrowPan(btn ArrowButton, step uint8) {
	switch btn {
	case ArrowLeft:
		p.PanValue = p.PanValue - uint16(step)
	case ArrowRight:
		p.PanValue = p.PanValue + uint16(step)
	}
}

func (p *PanTilt) HandleArrowTilt(btn ArrowButton, step uint8) {
	switch btn {
	case ArrowLeft:
		p.TiltValue = p.TiltValue - uint16(step)
	case ArrowRight:
		p.TiltValue = p.TiltValue + uint16(step)
	}
}
