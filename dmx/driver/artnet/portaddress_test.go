package artnet_test

import (
	"testing"

	"go.skymyer.dev/show-control/dmx/driver/artnet"
	"gotest.tools/assert"
)

func TestPortAddress(t *testing.T) {
	for _, test := range []struct {
		net, sub, uni uint8
		str           string
	}{
		{0, 0, 0, "0-0-0"},
		{1, 2, 3, "1-2-3"},
		{61, 15, 14, "61-15-14"},
		{127, 14, 15, "127-14-15"},
		{127, 15, 15, "127-15-15"},
	} {
		p, err := artnet.NewPortAddress(test.net, test.sub, test.uni)
		assert.NilError(t, err)

		assert.Equal(t, true, p.Valid())
		assert.Equal(t, test.str, p.String())
		assert.Equal(t, test.net, p.Net())
		assert.Equal(t, test.sub, p.SubNet())
		assert.Equal(t, test.uni, p.Universe())

		p, err = artnet.ParsePortAddress(test.str)
		assert.NilError(t, err)
		assert.Equal(t, test.str, p.String())
	}
}

func TestEmptyPortAddress(t *testing.T) {
	var zero uint8 = 0
	p := artnet.EmptyPortAddress()
	assert.Equal(t, true, p.Valid())
	assert.Equal(t, zero, p.Net())
	assert.Equal(t, zero, p.SubNet())
	assert.Equal(t, zero, p.Universe())
}

func TestInvalidPortAddress(t *testing.T) {
	for _, test := range []struct {
		net, sub, uni uint8
		err           error
	}{
		{128, 3, 14, artnet.ErrPortAddressNet},
		{53, 16, 14, artnet.ErrPortAddressSubNet},
		{53, 3, 16, artnet.ErrPortAddressUniverse},
	} {
		p, err := artnet.NewPortAddress(test.net, test.sub, test.uni)
		assert.Error(t, err, test.err.Error())
		assert.Equal(t, false, p.Valid())
		assert.Equal(t, artnet.InvalidPortAddress, p)
	}
}

func TestInvalidParsePortAddress(t *testing.T) {
	for _, test := range []struct {
		in  string
		err error
	}{
		{"", artnet.ErrPortAddressComponents},
		{"1", artnet.ErrPortAddressComponents},
		{"1-2", artnet.ErrPortAddressComponents},
		{"1-x", artnet.ErrPortAddressComponents},
		{"1-2-x", artnet.ErrPortAddressComponents},
		{"128-1-2", artnet.ErrPortAddressNet},
		{"127-16-2", artnet.ErrPortAddressSubNet},
		{"127-15-16", artnet.ErrPortAddressUniverse},
	} {
		p, err := artnet.ParsePortAddress(test.in)
		assert.Error(t, err, test.err.Error())
		assert.Equal(t, false, p.Valid())
		assert.Equal(t, artnet.InvalidPortAddress, p)
	}
}
