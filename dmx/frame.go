package dmx

import (
	"fmt"
	"sync"
)

var FrameLock sync.RWMutex

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
		slots: slots,
		data:  make([]byte, slots+1),
	}
	f.SetKind(kind)
	return f
}

type Frame struct {
	slots int
	data  []byte
}

func (f *Frame) SetKind(kind FrameKind) {
	f.data[0] = byte(kind)
}

func (f *Frame) SetSlot(slot int, val byte) error {
	// slot 0 is the start code
	if slot < 1 || slot > f.slots {
		return fmt.Errorf("invalid slot number")
	}
	f.data[slot] = val
	return nil
}

func (f *Frame) ClearAll() {
	kind := f.data[0]
	f.data = make([]byte, f.slots+1)
	f.data[0] = kind
}

func (f *Frame) Read() []byte {
	FrameLock.RLock()
	defer FrameLock.RUnlock()
	return f.data
}
