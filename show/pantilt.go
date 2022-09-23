package show

import (
	"sync"
	"time"

	"go.skymyer.dev/show-control/config"
)

func NewPantilt(spec *config.PantiltSpec) *Pantilt {
	return &Pantilt{
		fixtures: make(map[string]*PantiltFixture),
		spec:     spec,
		stopCh:   make(chan bool),
	}
}

type Pantilt struct {
	fixtures map[string]*PantiltFixture
	spec     *config.PantiltSpec
	stopCh   chan bool
	mutex    sync.Mutex
}

func (p *Pantilt) Start() error {

	ticker := time.NewTicker(time.Duration(100 * int(time.Millisecond)))

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-p.stopCh:
				return
			case <-ticker.C:
				continue
			}
		}
	}()
	return nil
}

func (p *Pantilt) Stop() error {
	p.stopCh <- true
	return nil
}

func (p *Pantilt) Apply() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for _, f := range p.fixtures {
		f.Set(byte(p.spec.Fixed.Pan), byte(p.spec.Fixed.Tilt))
	}
}

func (p *Pantilt) RegisterFixture(f *Fixture) error {
	cf := f.GetChanFuncByFeature(FEATURE_PAN, FEATURE_TILT)
	new := &PantiltFixture{
		fixture:   f,
		chanFuncs: cf,
	}
	p.fixtures[f.name] = new
	return nil
}

type PantiltFixture struct {
	fixture   *Fixture
	chanFuncs map[string][]string
}

func (p *PantiltFixture) Set(pan, tilt byte) {
	for _, cf := range p.chanFuncs[FEATURE_PAN] {
		ch, fn := SplitChanFunc(cf)
		p.fixture.MustSetValue(ch, fn, pan)
	}
	for _, cf := range p.chanFuncs[FEATURE_TILT] {
		ch, fn := SplitChanFunc(cf)
		p.fixture.MustSetValue(ch, fn, tilt)
	}
}
