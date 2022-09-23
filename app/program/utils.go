package program

import (
	"go.skymyer.dev/show-control/io"
)

func buttonWeightedUint8(options []io.Button, selected int) uint8 {
	return buttonWeightedForRangeUint8(options, selected, 0, ^uint8(0))
}

func buttonWeightedForRangeUint8(options []io.Button, selected int, min, max uint8) (value uint8) {
	return uint8((float64(max-min) / float64(len(options)-1)) * float64(selected))
}

func buttonWeightedUint16(options []io.Button, selected int) uint16 {
	return buttonWeightedForRangeUint16(options, selected, 0, ^uint16(0))
}

func buttonWeightedForRangeUint16(options []io.Button, selected int, min, max uint16) (value uint16) {
	return uint16((float64(max-min) / float64(len(options)-1)) * float64(selected))
}
