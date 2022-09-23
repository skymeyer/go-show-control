package program

import (
	"go.skymyer.dev/show-control/io"
)

type RGBW struct {
	Red   []io.Button
	Green []io.Button
	Blue  []io.Button
	White []io.Button

	RedValue   uint8
	GreenValue uint8
	BlueValue  uint8
	WhiteValue uint8
}

func (r *RGBW) Handle(button io.Button) bool {
	for index, btn := range r.Red {
		if btn == button {
			r.RedValue = buttonWeightedUint8(r.Red, index)
			return true
		}
	}
	for index, btn := range r.Green {
		if btn == button {
			r.GreenValue = buttonWeightedUint8(r.Green, index)
			return true
		}
	}
	for index, btn := range r.Blue {
		if btn == button {
			r.BlueValue = buttonWeightedUint8(r.Blue, index)
			return true
		}
	}
	for index, btn := range r.White {
		if btn == button {
			r.WhiteValue = buttonWeightedUint8(r.White, index)
			return true
		}
	}
	return false
}
