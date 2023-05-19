package show

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"go.skymyer.dev/show-control/app/logger"
	"go.skymyer.dev/show-control/common"
	"go.skymyer.dev/show-control/config"
	"go.skymyer.dev/show-control/io"
	"go.skymyer.dev/show-control/library/effect"
	"go.skymyer.dev/show-control/library/feature"
	"go.skymyer.dev/show-control/library/utils"
)

type ExecutorSequence struct {
	Name           string
	Duration       *time.Duration
	BeatsPerMinute *uint8
	Loop           *bool
}

func newExecutor(f map[string]*Fixture, g map[string]*config.Group) (*Executor, error) {
	return &Executor{
		fixtures: f,
		groups:   g,
	}, nil
}

type Executor struct {
	fixtures map[string]*Fixture
	groups   map[string]*config.Group

	show *Show
	lock sync.Mutex

	activeSequences []ExecutorSequence
}

func (e *Executor) Load(files ...string) error {
	var show = &Show{}
	for _, file := range files {
		if err := common.LoadFromFile(file, show); err != nil {
			return err
		}
	}

	// ensure sequence names are set
	for name, seq := range show.Sequences {
		seq.Name = name
	}

	e.show = show
	return nil
}

func (e *Executor) GetLiveSequences(page io.Page) []*Sequence {
	e.lock.Lock()
	defer e.lock.Unlock()

	// Filter out sequences with a button id configured
	var list []*Sequence
	for _, seq := range e.show.Sequences {
		if seq.Button > 0 && seq.Page == int8(page) {
			list = append(list, seq)
		}
	}

	return list
}

func (e *Executor) EnableSequences(list []ExecutorSequence) {
	e.lock.Lock()
	defer e.lock.Unlock()

	for _, es := range list {

		// Check if sequence is defined
		if _, ok := e.show.Sequences[es.Name]; !ok {
			logger.Default.Warn("enable sequence: unknown", zap.String("sequence", es.Name))
			continue
		}

		seq := e.show.Sequences[es.Name]

		// Override sequence settings if request
		if es.Duration != nil {
			seq.Duration = *es.Duration
		}
		if es.BeatsPerMinute != nil {
			seq.BeatsPerMinute = *es.BeatsPerMinute
		}
		if es.Loop != nil {
			seq.Loop = *es.Loop
		}

		// Initialize the sequence from scratch
		seq.Init()

		// Only add the sequence if it not yet in the active list
		onStack := false
		for i, active := range e.activeSequences {
			if active.Name == seq.Name {
				onStack = true
				logger.Default.Debug("enable sequences: already on stack", zap.String("sequence", seq.Name))
				e.activeSequences[i] = es
				continue
			}
		}
		if !onStack {
			e.activeSequences = append(e.activeSequences, es)
		}
	}

	e.garbageCollectSequences()
}

func (e *Executor) DisableSequences(names []string) {
	e.lock.Lock()
	defer e.lock.Unlock()

	for _, name := range names {
		e.activeSequences = removeExecutorSequence(e.activeSequences, name)
	}
	e.garbageCollectSequences()
}

func (e *Executor) GetSequenceStack() []ExecutorSequence {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.activeSequences
}

func (e *Executor) garbageCollectSequences() {
	for _, es := range e.activeSequences {
		if e.show.Sequences[es.Name].finished {
			logger.Default.Debug("gc remove sequence", zap.String("sequence", es.Name))
			e.activeSequences = removeExecutorSequence(e.activeSequences, es.Name)
		}
	}
	var stack []string
	for _, i := range e.activeSequences {
		stack = append(stack, i.Name)
	}
	logger.Default.Info("active sequences", zap.Int("count", len(e.activeSequences)), zap.Any("stack", stack))
}

func (e *Executor) renderFixtureFeatures(t time.Time) []Feature {
	e.lock.Lock()
	defer e.lock.Unlock()

	var (
		allCues    []Feature
		allEffects []Feature
	)

	for _, es := range e.activeSequences {
		cues, effects := e.show.Sequences[es.Name].GetCuesAndEffects(t)
		allCues = append(allCues, e.applyCues(cues)...)

		// Override BPM from sequence executor if set
		bpm := e.show.Sequences[es.Name].BeatsPerMinute
		if es.BeatsPerMinute != nil {
			bpm = *es.BeatsPerMinute
		}
		allEffects = append(allEffects, e.applyEffects(t, effects, bpm)...)
	}

	return append(allCues, allEffects...)
}

func (e *Executor) applyCues(cues []string) (features []Feature) {
	for _, cue := range cues {
		list, err := e.applyCue(cue)
		if err != nil {
			logger.Default.Warn("applyCues",
				zap.String("cue", cue),
				zap.Error(err),
			)
		}
		features = append(features, list...)
	}
	return features
}

func (e *Executor) applyCue(name string) (features []Feature, err error) {

	cueFeatures, ok := e.show.Cues[name]
	if !ok {
		return nil, fmt.Errorf("cue %q not found", name)
	}

	for _, feat := range cueFeatures {

		// Determine list of fixtures
		var fixtures = feat.Fixtures
		for _, g := range feat.Groups {
			if group, ok := e.groups[g]; ok {
				fixtures = append(fixtures, group.Members...)
			}
		}

		// Load feature config
		for _, f := range fixtures {
			fh, err := e.fixtures[f].GetFeature(feat.Feature)
			if err != nil {
				logger.Default.Warn("applyEffect invalid feature", zap.Error(err))
				continue
			}
			kind := fh.Kind()
			config := feature.NewConfig(kind, feat.Config)
			features = append(features, Feature{
				Name:    feat.Feature,
				Kind:    kind,
				Fixture: f,
				Config:  config,
			})
		}
	}
	return features, nil
}

func (e *Executor) applyEffects(t time.Time, effects []string, bpm uint8) (features []Feature) {
	for _, effect := range effects {
		list, err := e.applyEffect(t, effect, bpm)
		if err != nil {
			logger.Default.Warn("applyEffects",
				zap.String("effect", effect),
				zap.Error(err),
			)
		}
		features = append(features, list...)
	}
	return features
}

func (e *Executor) applyEffect(t time.Time, name string, bpm uint8) (features []Feature, err error) {

	fxs, ok := e.show.Effects[name]
	if !ok {
		return nil, fmt.Errorf("effect %q not found", name)
	}

	for _, fx := range fxs {

		fxHandler := effect.New(fx.Kind, fx.Config)
		fxHandler.Init(bpm)

		for _, modulator := range fx.Modulate {

			// If no explicit min/max is given, use full range by default
			if modulator.Min == nil {
				min := uint16(0)
				modulator.Min = &min
			}
			if modulator.Max == nil {
				max := ^uint16(0)
				modulator.Max = &max
			}

			amplitude := float64(*modulator.Max-*modulator.Min) / 2

			// Load & apply effect
			for id, collection := range modulator.Collections {

				// Render effect applying shift based on collection
				factor := fxHandler.Render(t, float64(id)*modulator.Shift)

				// Determine list of fixtures
				var fixtures = collection.Fixtures
				for _, g := range collection.Groups {
					if group, ok := e.groups[g]; ok {
						fixtures = append(fixtures, group.Members...)
					}
				}

				for _, f := range fixtures {

					fh, err := e.fixtures[f].GetFeature(modulator.Feature)
					if err != nil {
						logger.Default.Warn("applyEffect invalid feature", zap.Error(err))
						continue
					}

					modulated := uint16(amplitude+(factor*amplitude)) + *modulator.Min

					values := make(utils.FeatureValues)
					for _, property := range modulator.Properties {
						values[property] = modulated
					}

					features = append(features, Feature{
						Name:    modulator.Feature,
						Kind:    fh.Kind(),
						Fixture: f,
						Config:  fh.ConfigFromFeatureValues(values),
					})
				}
			}
		}
	}
	return features, nil
}

func removeExecutorSequence(in []ExecutorSequence, remove string) []ExecutorSequence {
	for i, v := range in {
		if v.Name == remove {
			return append(in[:i], in[i+1:]...)
		}
	}
	return in
}
