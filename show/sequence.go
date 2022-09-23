package show

import (
	"time"

	"go.uber.org/zap"

	"go.skymyer.dev/show-control/app/logger"
)

type SequenceStep struct {
	Cues     []string      `yaml:"cues,omitempty"`
	Effects  []string      `yaml:"fxs,omitempty"`
	Track    bool          `yaml:"track,omitempty"`
	Untrack  bool          `yaml:"untrack,omitempty"`
	Duration time.Duration `yaml:"duration,omitempty"`
}

type Sequence struct {
	Name           string         `yaml:"name,omitempty"`
	Steps          []SequenceStep `yaml:"steps,omitempty"`
	Duration       time.Duration  `yaml:"duration,omitempty"`
	BeatsPerMinute uint8          `yaml:"bpm,omitempty"`
	Loop           bool           `yaml:"loop,omitempty"`
	Manual         uint8          `yaml:"manual,omitempty"`

	started  bool
	finished bool

	start      time.Time
	transition time.Time
	step       int

	cuesActive     []string
	cuesTracked    []string
	effectsActive  []string
	effectsTracked []string
}

func (s *Sequence) Init() {
	// If no sequence duration is given, calculate it based on the steps.
	if s.Duration == 0 {
		for _, step := range s.Steps {
			s.Duration = s.Duration + step.Duration
		}
	}

	s.started = false
	s.finished = false
}

func (s *Sequence) Start(t time.Time) {
	s.start = t
	s.transition = t
	s.step = 0

	s.started = true
	s.finished = false

	s.cuesActive = s.Steps[0].Cues
	s.effectsActive = s.Steps[0].Effects

	if s.Steps[0].Track {
		s.cuesTracked = s.Steps[0].Cues
		s.effectsTracked = s.Steps[0].Effects
	} else {
		s.cuesTracked = []string{}
		s.effectsTracked = []string{}
	}

	logger.Default.Debug("sequence start",
		zap.String("sequence", s.Name),
		zap.Strings("cues", s.cuesActive),
		zap.Strings("effects", s.effectsActive),
		zap.Duration("duration", s.Duration),
		zap.Uint8("bpm", s.BeatsPerMinute),
	)
}

func (s *Sequence) GetCuesAndEffects(t time.Time) (cues []string, effects []string) {

	// Ensure the sequence has started
	if !s.started {
		s.Start(t)
	}

	// When marked as finished (non-loop), no cues/effects
	if s.finished {
		return cues, effects
	}

	// When sequence duration has expired, mark as finished.
	if s.start.Add(s.Duration).Before(t) {
		s.finished = true
		logger.Default.Debug("sequence finished", zap.String("sequence", s.Name), zap.String("reason", "sequence duration"))
		return cues, effects
	}

	step := s.Steps[s.step]

	// Switch to next step after duration of current step
	if s.transition.Add(step.Duration).Before(t) {

		s.transition = t
		s.step++

		if len(s.Steps) == s.step {
			if !s.Loop { // If no more steps and not looped, we finish the sequence
				s.finished = true
				logger.Default.Debug("sequence finished", zap.String("sequence", s.Name), zap.String("reason", "no more steps"))
				return cues, effects
			}
			logger.Default.Debug("sequence loop", zap.String("sequence", s.Name))
			s.step = 0
			s.cuesTracked = []string{}
			s.effectsTracked = []string{}
		}

		next := s.Steps[s.step]

		if next.Untrack {
			s.cuesTracked = []string{}
			s.effectsTracked = []string{}
		}

		s.cuesActive = append(s.cuesTracked, next.Cues...)
		s.effectsActive = append(s.effectsTracked, next.Effects...)

		if next.Track {
			s.cuesTracked = append(s.cuesTracked, next.Cues...)
			s.effectsTracked = append(s.effectsTracked, next.Effects...)
		}

		logger.Default.Debug("sequence next step",
			zap.String("sequence", s.Name),
			zap.Int("step", s.step),
			zap.Uint8("bpm", s.BeatsPerMinute),
			zap.Strings("cues", s.cuesActive),
			zap.Strings("effects", s.effectsActive),
		)
	}

	return s.cuesActive, s.effectsActive
}
