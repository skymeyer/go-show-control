package show

import (
	"go.skymyer.dev/show-control/config"
)

func LoadEffect(cfgs []*config.Effect, fixtures map[string]*Fixture, groups map[string]*config.Group) (*EffectCollection, error) {
	var e = &EffectCollection{}

	for _, cfg := range cfgs {

		// Instantiate effect
		var effect Effect
		switch cfg.Kind {
		case config.EFFECT_KIND_STROBE:
			spec, _ := cfg.Spec.(*config.StrobeSpec)
			effect = NewStrobe(spec)
		case config.EFFECT_KIND_MATRIX:
			spec, _ := cfg.Spec.(*config.MatrixSpec)
			effect = NewMatrix(spec)
		case config.EFFECT_KIND_PAN_TILT:
			spec, _ := cfg.Spec.(*config.PantiltSpec)
			effect = NewPantilt(spec)
		case config.EFFECT_KIND_RAW_FUNC:
			spec, _ := cfg.Spec.(*config.RawFunctionEffectSpec)
			effect = NewRawFunctionEffect(spec)
		case config.EFFECT_KIND_RAW_VALUE:
			spec, _ := cfg.Spec.(*config.RawValueEffectSpec)
			effect = NewRawValueEffect(spec)
		default:
			panic("unknown effect kind")
		}

		// Attach fixtures and register on collection
		for _, f := range GetFixtureCollection(fixtures, groups, cfg.Fixtures, cfg.Groups) {
			effect.RegisterFixture(f)
		}
		e.list = append(e.list, effect)
	}

	return e, nil
}

type EffectCollection struct {
	list []Effect
}

func (e *EffectCollection) Start() error {
	for _, efx := range e.list {
		efx.Start()
	}
	return nil
}

func (e *EffectCollection) Stop() error {
	for _, efx := range e.list {
		efx.Stop()
	}
	return nil
}

func (e *EffectCollection) Apply() {
	for _, efx := range e.list {
		efx.Apply()
	}
}

type Effect interface {
	RegisterFixture(*Fixture) error
	Start() error
	Stop() error
	Apply()
}
