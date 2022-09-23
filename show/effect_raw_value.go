package show

import (
	"sync"

	"go.skymyer.dev/show-control/config"
)

func NewRawValueEffect(spec *config.RawValueEffectSpec) *RawValueEffect {
	return &RawValueEffect{
		spec: spec,
	}
}

type RawValueEffect struct {
	fixtures []*Fixture
	spec     *config.RawValueEffectSpec
	mutex    sync.Mutex
	value    int
}

func (r *RawValueEffect) Start() error {
	r.value = r.spec.Value
	return nil
}

func (r *RawValueEffect) Stop() error {
	return nil
}

func (r *RawValueEffect) Apply() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, f := range r.fixtures {
		f.MustSetValue(r.spec.Channel, r.spec.Function, byte(r.value))
	}
}

func (r *RawValueEffect) RegisterFixture(f *Fixture) error {
	r.fixtures = append(r.fixtures, f)
	return nil
}
