package sine

import (
	"math"
	"time"
)

const (
	NAME = "sine"
)

type Effect struct {
	BeatsPerMinute uint8   `yaml:"bpm,omitempty"`
	BPMFraction    float64 `yaml:"bpmf,omitempty"`
	Clip           bool    `yaml:"clip,omitempty"`

	angularFreq float64
}

func (e *Effect) Init(bpm uint8) error {
	// Default BPM fraction to one
	if e.BPMFraction == 0 {
		e.BPMFraction = 1
	}
	// Override configured BPM from initialization if set
	if bpm > 0 {
		e.BeatsPerMinute = bpm
	}
	e.angularFreq = 2.0 * math.Pi * (float64(e.BeatsPerMinute) * e.BPMFraction / 60)
	return nil
}

func (e *Effect) Render(t time.Time, shift float64) float64 {
	// The shift defines the percentage of 360 degrees
	phase := 2.0 * math.Pi * shift

	sine := math.Sin(e.angularFreq*float64(t.UnixMilli())/1000 + phase)

	if e.Clip {
		if sine >= 0 {
			return 1
		} else {
			return -1
		}
	}

	return sine
}
