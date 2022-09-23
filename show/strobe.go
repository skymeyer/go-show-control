package show

import (
	"sync"

	"go.skymyer.dev/show-control/config"
)

func NewStrobe(spec *config.StrobeSpec) *Strobe {
	return &Strobe{
		spec: spec,
	}
}

type Strobe struct {
	fixtures []*StrobeFixture
	spec     *config.StrobeSpec
	mutex    sync.Mutex
	speed    byte
}

func (s *Strobe) Start() error {
	s.setSpeed(byte(s.spec.Speed))
	return nil
}

func (s *Strobe) Stop() error {
	s.setSpeed(byte(s.spec.Off))
	return nil
}

func (s *Strobe) Apply() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, f := range s.fixtures {
		f.SetSpeed(s.speed)
	}
}

func (s *Strobe) RegisterFixture(f *Fixture) error {
	cf := f.GetChanFuncByFeature(FEATURE_STROBE)
	new := &StrobeFixture{
		fixture:   f,
		chanFuncs: cf,
	}
	s.fixtures = append(s.fixtures, new)
	return nil
}

func (s *Strobe) setSpeed(speed byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.speed = speed
}

type StrobeFixture struct {
	fixture   *Fixture
	chanFuncs map[string][]string
}

func (s *StrobeFixture) SetSpeed(v byte) {
	for _, cf := range s.chanFuncs[FEATURE_STROBE] {
		ch, fn := SplitChanFunc(cf)
		s.fixture.MustSetValue(ch, fn, v)
	}
}
