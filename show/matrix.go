package show

import (
	"math"
	"math/bits"
	"sync"
	"time"

	"go.skymyer.dev/show-control/config"
)

const (
	MATRIX_EMPTY = "empty"
)

func NewMatrix(spec *config.MatrixSpec) *Matrix {
	return &Matrix{
		fixtures: make(map[string]*MatrixFixture),
		spec:     spec,
		stopCh:   make(chan bool),
	}
}

type Matrix struct {
	fixtures map[string]*MatrixFixture
	spec     *config.MatrixSpec
	stopCh   chan bool
	mutex    sync.Mutex
	mask     uint8
}

func (m *Matrix) Start() error {

	ticker := time.NewTicker(time.Duration(m.spec.Cycle * int(time.Millisecond)))
	m.mask = m.spec.Pattern

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-m.stopCh:
				m.setMask(255)
				return
			case <-ticker.C:
				m.setMask(bits.RotateLeft8(m.mask, m.spec.Step))
			}
		}
	}()
	return nil
}

func (m *Matrix) Stop() error {
	m.stopCh <- true
	return nil
}

func (m *Matrix) Apply() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for i, f := range m.spec.Matrix[0] {
		if f == MATRIX_EMPTY {
			continue
		}
		if int(math.Pow(2, float64(i)))&int(m.mask) == 0 {
			m.fixtures[f].SetDimmer(0x00)
		}
	}
}

func (m *Matrix) RegisterFixture(f *Fixture) error {
	cf := f.GetChanFuncByFeature(FEATURE_MASTER_DIMMER)
	new := &MatrixFixture{
		fixture:   f,
		chanFuncs: cf,
	}
	m.fixtures[f.name] = new
	return nil
}

func (m *Matrix) setMask(mask uint8) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.mask = mask
}

type MatrixFixture struct {
	fixture   *Fixture
	chanFuncs map[string][]string
}

func (m *MatrixFixture) SetDimmer(v byte) {
	for _, cf := range m.chanFuncs[FEATURE_MASTER_DIMMER] {
		ch, fn := SplitChanFunc(cf)
		m.fixture.MustSetValue(ch, fn, v)
	}
}
