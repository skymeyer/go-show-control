package dmx

import (
	"fmt"
	"sync"
)

// FIXME - where does this belong ?
type Channels map[string]*Channel

type Channel struct {
	Channel     int
	Default     int
	Description string
	Attributes  map[string]*Attribute
}

type Attribute struct {
	Min uint8
	Max uint8
}

type AttributeValue struct {
	Channel   string
	Attribute string
	Value     uint8
}

// DMX Frame Details and Frame Rate
//
// https://erg.abdn.ac.uk/users/gorry/eg3576/DMX-frame.html#:~:text=In%20DMX%2C%20the%20break%20(at,the%20start%20of%20a%20frame.

const (
	DMX512_SLOTS = 512

	FRAME_KIND_DEFAULT FrameKind = 0x00
	FRAME_KIND_TEXT    FrameKind = 0x17
	FRAME_KIND_RDM     FrameKind = 0xCC
	FRAME_KIND_SIP     FrameKind = 0xCF
)

type FrameKind byte

func NewDMX512Frame() *Frame {
	return NewFrame(FRAME_KIND_DEFAULT, 512)
}

func NewFrame(kind FrameKind, slots int) *Frame {
	f := &Frame{
		kind:    kind,
		slots:   slots,
		current: make([]byte, slots+1),
		new:     make([]byte, slots+1),
	}
	f.current[0] = byte(kind)
	f.new[0] = byte(kind)
	return f
}

type Frame struct {
	kind    FrameKind
	slots   int
	current []byte
	new     []byte
	apply   sync.Mutex
}

func (f *Frame) Apply(clear bool) {
	f.apply.Lock()
	defer f.apply.Unlock()

	copy(f.current, f.new)

	if clear {
		f.new = make([]byte, f.slots+1)
		f.new[0] = byte(f.kind)
	}
}

func (f *Frame) SetSlot(slot int, val byte) error {
	if slot < 1 || slot > f.slots {
		return fmt.Errorf("invalid slot number")
	}
	f.new[slot] = val
	return nil
}

func (f *Frame) GetSlot(slot int) (byte, error) {
	if slot < 1 || slot > f.slots {
		return 0, fmt.Errorf("invalid slot number")
	}
	return f.new[slot], nil
}

func (f *Frame) Read() []byte {
	f.apply.Lock()
	defer f.apply.Unlock()
	return f.current
}
