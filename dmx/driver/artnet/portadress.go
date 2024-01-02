package artnet

import (
	"fmt"
)

const (
	InvalidPortAddress PortAddress = 0b1000000000000000 // Invalid port address with first bit set to one
	PortAddressSep                 = "-"
)

var (
	ErrPortAddressNet        = fmt.Errorf("port address net out of range [0-127]")
	ErrPortAddressSubNet     = fmt.Errorf("port address subnet out of range [0-15]")
	ErrPortAddressUniverse   = fmt.Errorf("port address universe out of range [0-15]")
	ErrPortAddressComponents = fmt.Errorf("port address requires 3 components 1-2-3")
)

type PortAddress uint16

func ParsePortAddress(in string) (PortAddress, error) {
	var (
		net      uint8
		subnet   uint8
		universe uint8
	)
	n, err := fmt.Sscanf(in, "%d-%d-%d", &net, &subnet, &universe)
	if err != nil || n != 3 {
		return InvalidPortAddress, ErrPortAddressComponents
	}
	return NewPortAddress(net, subnet, universe)
}

func NewPortAddress(net, subnet, universe uint8) (PortAddress, error) {
	if net > 127 {
		return InvalidPortAddress, ErrPortAddressNet
	}
	if subnet > 15 {
		return InvalidPortAddress, ErrPortAddressSubNet
	}
	if universe > 15 {
		return InvalidPortAddress, ErrPortAddressUniverse
	}
	return PortAddress(uint16(net)<<8 | uint16(subnet)<<4 | uint16(universe)), nil
}

func MustPortAddress(net, subnet, universe uint8) PortAddress {
	pa, err := NewPortAddress(net, subnet, universe)
	if err != nil {
		panic(fmt.Sprintf("port address: %v", err))
	}
	return pa
}

func EmptyPortAddress() PortAddress {
	return MustPortAddress(0, 0, 0)
}

func (p PortAddress) Valid() bool {
	return p&0b1000000000000000 == 0
}

func (p PortAddress) Net() uint8 {
	return uint8(p >> 8)
}

func (p PortAddress) SubNet() uint8 {
	return uint8(p & 0b0000000011110000 >> 4)
}

func (p PortAddress) Universe() uint8 {
	return uint8(p & 0b0000000000001111)
}

func (p PortAddress) String() string {
	return fmt.Sprintf("%d-%d-%d", p.Net(), p.SubNet(), p.Universe())
}
