package artnet_test

import (
	"testing"

	"gotest.tools/assert"

	"go.skymyer.dev/show-control/dmx/driver/artnet"
)

func TestTalkToMe(t *testing.T) {
	for _, test := range []struct {
		f artnet.TalkToMeFlag
		r uint8
	}{
		{
			f: artnet.TTMReplyOnChange,
			r: 0b00000010,
		},
		{
			f: artnet.TTMDiagnostics,
			r: 0b00000100,
		},
		{
			f: artnet.TTMDiagnosticsUnicast,
			r: 0b00001000,
		},
		{
			f: artnet.TTMDisableVLC,
			r: 0b00010000,
		},
		{
			f: artnet.TTMEnableTargetedMode,
			r: 0b00100000,
		},
	} {
		nb := artnet.NewTalkToMe().With(test.f)
		assert.Equal(t, test.r, uint8(nb))
		assert.Equal(t, true, nb.Has(test.f))
	}
}

func TestTalkToMeWithWithout(t *testing.T) {
	for _, test := range []struct {
		with    []artnet.TalkToMeFlag
		without []artnet.TalkToMeFlag
		has     []artnet.TalkToMeFlag
		hasnot  []artnet.TalkToMeFlag
	}{
		{
			with: []artnet.TalkToMeFlag{
				artnet.TTMReplyOnChange,
				artnet.TTMDiagnostics,
				artnet.TTMDiagnosticsUnicast,
				artnet.TTMDisableVLC,
				artnet.TTMEnableTargetedMode,
			},
			without: []artnet.TalkToMeFlag{},
			has: []artnet.TalkToMeFlag{
				artnet.TTMReplyOnChange,
				artnet.TTMDiagnostics,
				artnet.TTMDiagnosticsUnicast,
				artnet.TTMDisableVLC,
				artnet.TTMEnableTargetedMode,
			},
			hasnot: []artnet.TalkToMeFlag{},
		},
		{
			with: []artnet.TalkToMeFlag{
				artnet.TTMReplyOnChange,
				artnet.TTMDiagnostics,
				artnet.TTMDiagnosticsUnicast,
				artnet.TTMDisableVLC,
				artnet.TTMEnableTargetedMode,
			},
			without: []artnet.TalkToMeFlag{
				artnet.TTMDiagnosticsUnicast,
				artnet.TTMDisableVLC,
			},
			has: []artnet.TalkToMeFlag{
				artnet.TTMReplyOnChange,
				artnet.TTMDiagnostics,
				artnet.TTMEnableTargetedMode,
			},
			hasnot: []artnet.TalkToMeFlag{
				artnet.TTMDiagnosticsUnicast,
				artnet.TTMDisableVLC,
			},
		},
	} {
		nb := artnet.NewTalkToMe()
		for _, opt := range test.with {
			nb = nb.With(opt)
		}
		for _, opt := range test.without {
			nb = nb.Without(opt)
		}
		for _, f := range test.has {
			assert.Equal(t, true, nb.Has(f))
		}
		for _, f := range test.hasnot {
			assert.Equal(t, false, nb.Has(f))
		}
	}
}
