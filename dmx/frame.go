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

// Frame maintains two byte registers, being `current` and `new`. Any values stored
// using `SetSlot` or read using `GetSlot` are using the `new` register. This allows
// logic to set/get values separately from logic sending the full DMX frame to a DMX
// device. Once all values of a given frame are in the correct state, calling `Apply`
// copies the state from `new` to `current`.
type Frame struct {
	kind    FrameKind
	slots   int
	current []byte
	new     []byte
	apply   sync.Mutex
}

// GetKind returns the type of DMX frame.
func (f *Frame) GetKind() FrameKind {
	return f.kind
}

// Apply the state of the `new` to the `current` register ready for Frame.Read.
func (f *Frame) Apply(clear bool) {
	f.apply.Lock()
	defer f.apply.Unlock()

	copy(f.current, f.new)

	if clear {
		f.new = make([]byte, f.slots+1)
		f.new[0] = byte(f.kind)
	}
}

// SetSlot sets the byte value of a given slot in the `new` register.
func (f *Frame) SetSlot(slot int, val byte) error {
	if slot < 1 || slot > f.slots {
		return fmt.Errorf("invalid slot number")
	}
	f.new[slot] = val
	return nil
}

// GetSlot returns the byte value of a given slot from the `new` register.
func (f *Frame) GetSlot(slot int) (byte, error) {
	if slot < 1 || slot > f.slots {
		return 0, fmt.Errorf("invalid slot number")
	}
	return f.new[slot], nil
}

// Read returns the full DMX frame including including kind from `current` register.
func (f *Frame) Read() []byte {
	f.apply.Lock()
	defer f.apply.Unlock()
	return f.current
}

// ReadSlots returns only the DMX slot data without the kind from `current` register.
func (f *Frame) ReadSlots() (data []byte) {
	f.apply.Lock()
	defer f.apply.Unlock()
	copy(data[:], f.current[1:])
	return data
}
