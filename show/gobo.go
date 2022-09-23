package show

import "go.skymyer.dev/show-control/config"

type Gobo struct {
	fixtures []*GoboFixture
	spec     *config.GoboSpec
}

func (gobo *Gobo) Apply() {
	gobo.Set(gobo.spec.Gobo)
}

func (gobo *Gobo) RegisterFixture(f *Fixture) error {
	cf := f.GetChanFuncByFeature(FEATURE_GOBO)
	new := &GoboFixture{
		fixture:   f,
		chanFuncs: cf,
	}
	gobo.fixtures = append(gobo.fixtures, new)
	return nil
}

func (gobo *Gobo) Set(v string) {
	for _, f := range gobo.fixtures {
		f.Set(v)
	}
}

type GoboFixture struct {
	fixture   *Fixture
	chanFuncs map[string][]string
}

func (gobo *GoboFixture) Set(v string) {
	for _, cf := range gobo.chanFuncs[FEATURE_GOBO] {
		ch, _ := SplitChanFunc(cf)
		gobo.fixture.MustSelectFunction(ch, v)
	}
}
