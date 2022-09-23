package show

import (
	"sync"

	"go.skymyer.dev/show-control/config"
)

func NewRawFunctionEffect(spec *config.RawFunctionEffectSpec) *RawFunctionEffect {
	return &RawFunctionEffect{
		spec: spec,
	}
}

type RawFunctionEffect struct {
	fixtures []*Fixture
	spec     *config.RawFunctionEffectSpec
	mutex    sync.Mutex
	fn       string
}

func (r *RawFunctionEffect) Start() error {
	r.fn = r.spec.FunctionStart
	return nil
}

func (r *RawFunctionEffect) Stop() error {
	r.fn = r.spec.FunctionStop
	return nil
}

func (r *RawFunctionEffect) Apply() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, f := range r.fixtures {
		f.MustSelectFunction(r.spec.Channel, r.fn)
	}
}

func (r *RawFunctionEffect) RegisterFixture(f *Fixture) error {
	r.fixtures = append(r.fixtures, f)
	return nil
}
