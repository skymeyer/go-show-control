package program

import (
	"go.skymyer.dev/show-control/io"
)

type Slider struct {
	Buttons []io.Button
	Value   uint8
}

func (s *Slider) Handle(button io.Button) bool {
	for index, btn := range s.Buttons {
		if btn == button {
			s.Value = buttonWeightedUint8(s.Buttons, index)
			return true
		}
	}
	return false
}

func (s *Slider) HandleArrow(btn ArrowButton, step uint8) {
	switch btn {
	case ArrowLeft:
		if s.Value < step {
			s.Value = 0
		} else {
			s.Value = s.Value - step
		}
	case ArrowRight:
		if s.Value > ^uint8(0)-step {
			s.Value = ^uint8(0)
		} else {
			s.Value = s.Value + step
		}
	}
}
