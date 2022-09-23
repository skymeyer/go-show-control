package show

import (
	"fmt"
	"strings"

	"go.skymyer.dev/show-control/config"
	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/library"
	"go.skymyer.dev/show-control/utils"
)

const (
	FEATURE_MASTER_DIMMER = "master_dimmer"
	FEATURE_RGB_RED       = "rgb_red"
	FEATURE_RGB_GREEN     = "rgb_green"
	FEATURE_RGB_BLUE      = "rgb_blue"
	FEATURE_COLOR         = "color"
	FEATURE_GOBO          = "gobo"
	FEATURE_STROBE        = "strobe"
	FEATURE_PAN           = "pan"
	FEATURE_TILT          = "tilt"
)

func NewFixture(lib library.Fixtures, name string, conf *config.Fixture, out *dmx.Frame) (*Fixture, error) {
	if _, ok := lib[conf.Kind]; !ok {
		return nil, fmt.Errorf("fixture kind %q not found", conf.Kind)
	}
	if _, ok := lib[conf.Kind].Modes[conf.Mode]; !ok {
		return nil, fmt.Errorf("invalid mode %q for %q", conf.Mode, conf.Kind)
	}
	f := &Fixture{
		name:     name,
		config:   conf,
		channels: lib[conf.Kind].Modes[conf.Mode],
		output:   out,
	}
	return f, nil
}

type Fixture struct {
	name     string
	config   *config.Fixture
	channels map[string]*library.Channel
	output   *dmx.Frame
}

func (f *Fixture) MustSetValue(channel, function string, val byte) {
	if err := f.SetValue(channel, function, val); err != nil {
		panic(err)
	}
}

func (f *Fixture) SetValue(channel, function string, val byte) error {
	if err := f.validateChannelAndFunc(channel, function); err != nil {
		return err
	}

	if val > byte(f.channels[channel].Functions[function].Max) || val < byte(f.channels[channel].Functions[function].Min) {
		return fmt.Errorf("value %d out of range for %s.%s", val, channel, function)
	}

	slot := f.channels[channel].Channel + f.config.Address - 1
	f.output.SetSlot(slot, val)
	return nil
}

func (f *Fixture) MustSelectFunction(channel, function string) {
	if err := f.SelectFunction(channel, function); err != nil {
		panic(err)
	}
}

func (f *Fixture) SelectFunction(channel, function string) error {
	if err := f.validateChannelAndFunc(channel, function); err != nil {
		return err
	}

	slot := f.channels[channel].Channel + f.config.Address - 1
	f.output.SetSlot(slot, byte(f.channels[channel].Functions[function].Min))
	return nil
}

func (f *Fixture) validateChannelAndFunc(channel, function string) error {
	if _, ok := f.channels[channel]; !ok {
		return fmt.Errorf("invalid channel %q", channel)
	}

	if _, ok := f.channels[channel].Functions[function]; !ok {
		return fmt.Errorf("invalid channel function %q", function)
	}
	return nil
}

func (f *Fixture) GetChanFuncByFeature(features ...string) map[string][]string {
	var result = make(map[string][]string)
	for _, feature := range features {
		for chName, ch := range f.channels {
			for fnName, fn := range ch.Functions {
				if utils.StringsContains(fn.Features, feature) {
					result[feature] = append(result[feature], JoinChanFunc(chName, fnName))
				}
			}
		}
	}
	return result
}

const (
	CHAN_FUNC_SEP = ":"
)

func JoinChanFunc(channel, function string) string {
	return fmt.Sprintf("%s%s%s", channel, CHAN_FUNC_SEP, function)
}

func SplitChanFunc(in string) (string, string) {
	res := strings.Split(in, CHAN_FUNC_SEP)
	return res[0], res[1]
}
