package show

import (
	"go.skymyer.dev/show-control/config"
)

func NewScene(cfg []*config.Scene, fixtures map[string]*Fixture, groups map[string]*config.Group) (*Scene, error) {
	s := &Scene{}

	for _, scene := range cfg {
		switch scene.Kind {
		case config.SCENE_KIND_RGB:
			spec, _ := scene.Spec.(*config.RGBSpec)
			new := &RGB{spec: spec}
			for _, f := range GetFixtureCollection(fixtures, groups, scene.Fixtures, scene.Groups) {
				new.RegisterFixture(f)
			}
			s.rgbs = append(s.rgbs, new)
		case config.SCENE_KIND_COLOR:
			spec, _ := scene.Spec.(*config.ColorSpec)
			new := &Color{spec: spec}
			for _, f := range GetFixtureCollection(fixtures, groups, scene.Fixtures, scene.Groups) {
				new.RegisterFixture(f)
			}
			s.colors = append(s.colors, new)
		case config.SCENE_KIND_GOBO:
			spec, _ := scene.Spec.(*config.GoboSpec)
			new := &Gobo{spec: spec}
			for _, f := range GetFixtureCollection(fixtures, groups, scene.Fixtures, scene.Groups) {
				new.RegisterFixture(f)
			}
			s.gobos = append(s.gobos, new)
		default:
			panic("unknown scene kind")
		}
	}

	return s, nil
}

type Scene struct {
	rgbs   []*RGB
	colors []*Color
	gobos  []*Gobo
}

func (s *Scene) Apply() error {
	for _, i := range s.rgbs {
		i.Apply()
	}
	for _, i := range s.colors {
		i.Apply()
	}
	for _, i := range s.gobos {
		i.Apply()
	}
	return nil
}
