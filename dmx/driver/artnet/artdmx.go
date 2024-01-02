package artnet

import (
	"bytes"
	"encoding/binary"

	"go.skymyer.dev/show-control/common"
)

type ArtDmx struct {
	Sequence uint8
	Physical uint8
	SubUni   uint8
	Net      uint8
	LengthHi uint8
	LengthLo uint8
	Data     [512]byte
}

func (p *ArtDmx) SetPortAddress(a PortAddress) {
	p.SubUni = uint8(a)
	p.Net = a.Net()
}

func (p *ArtDmx) SetData(data [512]byte) {
	p.Data = data
	p.LengthHi = 0x02
	p.LengthLo = 0x00
}

func (p *ArtDmx) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(common.PadToSize(data, p))
	return binary.Read(buf, binary.BigEndian, p)
}

func (p *ArtDmx) MarshalBinary() ([]byte, error) {
	buf, err := NewArtNetHeader(OpDmx).MarshalBuffer()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
