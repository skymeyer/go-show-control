package artnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	ArtNetVersion uint8 = 14
)

var (
	ArtNetID = [8]byte{0x41, 0x72, 0x74, 0x2d, 0x4e, 0x65, 0x74, 0x00}
)

var (
	ErrHeaderLength    = fmt.Errorf("invalid header length")
	ErrHeaderUnmarshal = fmt.Errorf("unmarshal header")
	ErrHeaderID        = fmt.Errorf("invalid artnet id")
	ErrHeaderVersion   = fmt.Errorf("protocol version too low")
)

func NewArtNetHeader(op OpCode) *ArtNetHeader {
	return &ArtNetHeader{
		ID:        ArtNetID,
		OpCode:    op,
		ProtVerLo: ArtNetVersion,
	}
}

type ArtNetHeader struct {
	ID        [8]byte
	OpCode    OpCode
	ProtVerHi uint8
	ProtVerLo uint8
}

func (h *ArtNetHeader) UnmarshalBinary(data []byte) error {
	if len(data) < 12 {
		return ErrHeaderLength
	}
	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.LittleEndian, h); err != nil {
		return ErrHeaderUnmarshal
	}

	if h.ID != ArtNetID {
		return ErrHeaderID
	}

	// ArtPollReply does not contain a protocol version field
	if h.OpCode == OpPollReply {
		h.ProtVerHi = 0
		h.ProtVerLo = 0
	} else {
		if h.ProtVerLo < ArtNetVersion {
			return ErrHeaderVersion
		}
	}

	return nil
}

func (h *ArtNetHeader) MarshalBuffer() (bytes.Buffer, error) {

	var buf bytes.Buffer

	if h.OpCode == 0 {
		return buf, fmt.Errorf("opcode not set")
	}

	// Ensure ID and protocol version is set
	h.ID = ArtNetID
	h.ProtVerLo = ArtNetVersion

	if err := binary.Write(&buf, binary.LittleEndian, h); err != nil {
		return buf, err
	}

	// ArtPollReply does not carry any version fields
	if h.OpCode == OpPollReply {
		buf.Truncate(10)
	}

	return buf, nil
}
