package show

import (
	"fmt"

	"encoding/hex"

	"go.skymyer.dev/show-control/config"
)

type RGB struct {
	fixtures []*RGBFixture
	spec     *config.RGBSpec
}

func (rgb *RGB) Apply() {
	rgb.SetColor(rgb.spec.Color)
}

func (rgb *RGB) RegisterFixture(f *Fixture) error {
	cf := f.GetChanFuncByFeature(FEATURE_RGB_RED, FEATURE_RGB_GREEN, FEATURE_RGB_BLUE)
	if len(cf) != 3 {
		return fmt.Errorf("missing rgb features")
	}
	new := &RGBFixture{
		fixture:   f,
		chanFuncs: cf,
	}
	rgb.fixtures = append(rgb.fixtures, new)
	return nil
}

func (rgb *RGB) SetColor(in string) error {
	v, err := hex.DecodeString(in)
	if err != nil {
		return fmt.Errorf("invalid color value %q: %v", in, err)
	}
	for _, f := range rgb.fixtures {
		f.SetColor8Bit(v)
	}
	return nil
}

type RGBFixture struct {
	fixture   *Fixture
	chanFuncs map[string][]string
}

func (rgb *RGBFixture) SetColor8Bit(c []byte) {
	for _, cf := range rgb.chanFuncs[FEATURE_RGB_RED] {
		ch, fn := SplitChanFunc(cf)
		rgb.fixture.MustSetValue(ch, fn, c[0])
	}
	for _, cf := range rgb.chanFuncs[FEATURE_RGB_GREEN] {
		ch, fn := SplitChanFunc(cf)
		rgb.fixture.MustSetValue(ch, fn, c[1])
	}
	for _, cf := range rgb.chanFuncs[FEATURE_RGB_BLUE] {
		ch, fn := SplitChanFunc(cf)
		rgb.fixture.MustSetValue(ch, fn, c[2])
	}
}
