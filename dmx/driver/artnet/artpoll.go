package artnet

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"go.skymyer.dev/show-control/common"
)

type ArtPoll struct {
	Flags                     TalkToMe
	DiagPriority              DpCode
	TargetPortAddressTopHi    uint8
	TargetPortAddressTopLo    uint8
	TargetPortAddressBottomHi uint8
	TargetPortAddressBottomLo uint8
	EstaManHi                 uint8
	EstaManLo                 uint8
	OemHi                     uint8
	OemLo                     uint8
}

func (p *ArtPoll) SetTargetPortAddress(top, bottom PortAddress) {
	p.TargetPortAddressTopHi, p.TargetPortAddressTopLo = common.Uint16ToUint8(uint16(top))
	p.TargetPortAddressBottomHi, p.TargetPortAddressBottomLo = common.Uint16ToUint8(uint16(bottom))
}

func (p *ArtPoll) GetTargetPortAddress() (top, bottom PortAddress) {
	top = PortAddress(common.Uint8ToUint16(p.TargetPortAddressTopHi, p.TargetPortAddressTopLo))
	bottom = PortAddress(common.Uint8ToUint16(p.TargetPortAddressBottomHi, p.TargetPortAddressBottomLo))
	return top, bottom
}

func (p *ArtPoll) SetESTAManCode(c ESTAManCode) {
	p.SetESTAManCodeRaw(uint16(c))
}

func (p *ArtPoll) GetESTAManCode() ESTAManCode {
	return ESTAManCode(p.GetESTAManCodeRaw())
}

func (p *ArtPoll) SetESTAManCodeRaw(c uint16) {
	p.EstaManHi, p.EstaManLo = common.Uint16ToUint8(c)
}

func (p *ArtPoll) GetESTAManCodeRaw() uint16 {
	return common.Uint8ToUint16(p.EstaManHi, p.EstaManLo)
}

func (p *ArtPoll) SetOemCode(c OemCode) {
	p.SetOemCodeRaw(uint16(c))
}

func (p *ArtPoll) GetOemCode() OemCode {
	return OemCode(p.GetOemCodeRaw())
}

func (p *ArtPoll) SetOemCodeRaw(c uint16) {
	p.OemHi, p.OemLo = common.Uint16ToUint8(c)
}

func (p *ArtPoll) GetOemCodeRaw() uint16 {
	return common.Uint8ToUint16(p.OemHi, p.OemLo)
}

func (p *ArtPoll) UnmarshalBinary(data []byte) error {
	if len(data) < 2 {
		return fmt.Errorf("ArtPoll expects at least 2 bytes of data, got %d bytes", len(data))
	}
	buf := bytes.NewReader(common.PadToSize(data, p))
	return binary.Read(buf, binary.LittleEndian, p)
}

func (p *ArtPoll) MarshalBinary() ([]byte, error) {
	buf, err := NewArtNetHeader(OpPoll).MarshalBuffer()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
