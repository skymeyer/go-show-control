package show

import "go.skymyer.dev/show-control/config"

type Color struct {
	fixtures []*ColorFixture
	spec     *config.ColorSpec
}

func (c *Color) Apply() {
	c.Set(c.spec.Color)
}

func (c *Color) RegisterFixture(f *Fixture) error {
	cf := f.GetChanFuncByFeature(FEATURE_COLOR)
	new := &ColorFixture{
		fixture:   f,
		chanFuncs: cf,
	}
	c.fixtures = append(c.fixtures, new)
	return nil
}

func (c *Color) Set(v string) {
	for _, f := range c.fixtures {
		f.Set(v)
	}
}

type ColorFixture struct {
	fixture   *Fixture
	chanFuncs map[string][]string
}

func (c *ColorFixture) Set(v string) {
	for _, cf := range c.chanFuncs[FEATURE_COLOR] {
		ch, _ := SplitChanFunc(cf)
		c.fixture.MustSelectFunction(ch, v)
	}
}
