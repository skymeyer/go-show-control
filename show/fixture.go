package show

import (
	"fmt"

	"go.skymyer.dev/show-control/config"
	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/library"
	"go.skymyer.dev/show-control/library/feature"
)

func NewFixture(lib library.Fixtures, name string, conf *config.Fixture, out *dmx.Frame) (*Fixture, error) {
	if _, ok := lib[conf.Kind]; !ok {
		return nil, fmt.Errorf("fixture kind %q not found", conf.Kind)
	}
	if _, ok := lib[conf.Kind].Modes[conf.Mode]; !ok {
		return nil, fmt.Errorf("invalid mode %q for %q", conf.Mode, conf.Kind)
	}

	var channels = lib[conf.Kind].Modes[conf.Mode].Dmx

	features := make(map[string]feature.Handler)
	for n, f := range lib[conf.Kind].Modes[conf.Mode].Features {
		features[n] = feature.NewHandler(f.Kind, f.Spec, channels)
	}

	f := &Fixture{
		name:     name,
		config:   conf,
		channels: channels,
		features: features,
		output:   out,
	}
	return f, nil
}

func NewVirtualFixture(lib library.Fixtures, name string, conf, realConf *config.Fixture, real *Fixture) (*Fixture, error) {
	if conf.Kind != "virtual" {
		return nil, fmt.Errorf("virtual device %q does not have virtual kind (have %q)", name, conf.Kind)
	}

	var channels = lib[realConf.Kind].Modes[realConf.Mode].Dmx

	features := make(map[string]feature.Handler)
	for n, f := range lib[realConf.Kind].Modes[realConf.Mode].Features {
		if mapped, ok := conf.Map[n]; ok {
			features[mapped] = feature.NewHandler(f.Kind, f.Spec, channels)
		}
	}

	f := &Fixture{
		name:     name,
		config:   real.config,
		channels: channels,
		features: features,
		output:   real.output,
	}
	return f, nil
}

type Fixture struct {
	name     string
	config   *config.Fixture
	channels dmx.Channels
	features map[string]feature.Handler
	output   *dmx.Frame
}

func (f *Fixture) GetFeature(name string) (feature.Handler, error) {
	if h, ok := f.features[name]; ok {
		return h, nil
	}
	return nil, fmt.Errorf("unknown feature %q for fixture %q", name, f.name)
}

func (f *Fixture) SetFeature(name string, config interface{}) (errs []error) {
	ats, err := f.features[name].Render(config)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	for _, at := range ats {
		if err := f.SetFeatureValue(at.Channel, at.Attribute, at.Value); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// SetFeatureValue converts a relative value (0-255) and transform it to
// the range as specified by the attribute. This allows for consistent
// interaction through features instead of dealing with specific (sub)
// ranges of the raw DMX values.
func (f *Fixture) SetFeatureValue(ch, attr string, val uint8) error {
	if err := f.validateChannelAttribute(ch, attr); err != nil {
		return err
	}

	var (
		min = f.channels[ch].Attributes[attr].Min
		max = f.channels[ch].Attributes[attr].Max
	)
	f.setValue(ch, uint8(float64(val)*(float64(max-min)/float64(^uint8(0))))+min)
	return nil
}

func (f *Fixture) GetFeatureValue(ch, attr string, val uint8) error {
	if err := f.validateChannelAttribute(ch, attr); err != nil {
		return err
	}

	// FIXME convert back to relative values and do something with attribute

	f.getValue(ch)
	return nil
}

// Set explicit DMX channel/attribute value as is.
func (f *Fixture) SetValue(channel, attr string, val byte) error {
	if err := f.validateChannelAttribute(channel, attr); err != nil {
		return err
	}

	if val > byte(f.channels[channel].Attributes[attr].Max) || val < byte(f.channels[channel].Attributes[attr].Min) {
		return fmt.Errorf("value %d out of range for %s.%s", val, channel, attr)
	}
	f.setValue(channel, val)
	return nil
}

func (f *Fixture) GetAttributeValue(channel string) (attr string, val byte, err error) {
	if err := f.validateChannel(channel); err != nil {
		return "", 0, err
	}
	val = f.getValue(channel)
	for attr, a := range f.channels[channel].Attributes {
		if val >= a.Min && val <= a.Max {
			return attr, val, nil
		}
	}
	return "", 0, fmt.Errorf("no valid attribute found on fixture %q for channel %q", f.name, channel)
}

func (f *Fixture) GetValue(channel string) (byte, error) {
	if err := f.validateChannel(channel); err != nil {
		return 0, err
	}
	return f.getValue(channel), nil
}

func (f *Fixture) setValue(channel string, val byte) {
	slot := f.channels[channel].Channel + f.config.Address - 1
	f.output.SetSlot(slot, val)
}

func (f *Fixture) getValue(channel string) byte {
	slot := f.channels[channel].Channel + f.config.Address - 1
	val, _ := f.output.GetSlot(slot)
	return val
}

func (f *Fixture) SetAttribute(channel, attr string) error {
	if err := f.validateChannelAttribute(channel, attr); err != nil {
		return err
	}

	slot := f.channels[channel].Channel + f.config.Address - 1
	f.output.SetSlot(slot, byte(f.channels[channel].Attributes[attr].Min))
	return nil
}

func (f *Fixture) validateChannel(ch string) error {
	if _, ok := f.channels[ch]; !ok {
		return fmt.Errorf("invalid channel %q", ch)
	}
	return nil
}

func (f *Fixture) validateChannelAttribute(ch, attr string) error {
	if err := f.validateChannel(ch); err != nil {
		return err
	}
	if _, ok := f.channels[ch].Attributes[attr]; !ok {
		return fmt.Errorf("invalid channel attribute %q for channel %q", attr, ch)
	}
	return nil
}
