package artnet

import (
	"encoding"
	"fmt"
)

const (
	ArtNetPort = 6454
)

type ArtNetPacket interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

func UnmarshalBinary(data []byte) (*ArtNetHeader, ArtNetPacket, error) {
	var (
		h = &ArtNetHeader{}
		p ArtNetPacket
	)
	if err := h.UnmarshalBinary(data); err != nil {
		return nil, nil, err
	}

	switch h.OpCode {
	case OpPoll:
		p = &ArtPoll{}
	case OpPollReply:
		p = &ArtPollReply{}
	default:
		return nil, nil, fmt.Errorf("not handling opcode %#v", h.OpCode)
	}

	// ArtPollReply does not carry version fields
	next := data[12:]
	if h.OpCode == OpPollReply {
		next = data[10:]
	}

	if err := p.UnmarshalBinary(next); err != nil {
		return nil, nil, err
	}

	return h, p, nil
}
