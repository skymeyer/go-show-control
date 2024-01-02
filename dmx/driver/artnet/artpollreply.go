package artnet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"go.skymyer.dev/show-control/app/logger"
	"go.skymyer.dev/show-control/common"
	"go.uber.org/zap"
)

type ArtPollReply struct {
	IPAddress     [4]byte
	Port          uint16
	VersInfoH     uint8
	VersInfoL     uint8
	NetSwitch     uint8
	SubSwitch     uint8
	OemHi         uint8
	OemLo         uint8
	UBEAVersion   uint8
	Status1       uint8
	EstaManLo     uint8
	EstaManHi     uint8
	PortName      [18]byte
	LongName      [64]byte
	NodeReport    [64]byte
	NumPortsHi    uint8
	NumPortsLo    uint8
	PortTypes     [4]byte
	GoodInput     [4]byte
	GoodOutputA   [4]byte
	SwIn          [4]byte
	SwOut         [4]byte
	AcnPriority   uint8
	SwMacro       uint8
	SwRemote      uint8
	_             [3]byte
	Style         StyleCode
	MACH          [6]byte
	BindIP        [4]byte
	BindIndex     uint8
	Status2       uint8
	GoodOutputB   [4]byte
	Status3       uint8
	DefaulRespUID [6]byte
	UserHi        uint8
	UserLo        uint8
	RefreshRateHi uint8
	RefreshRateLo uint8
	_             [11]byte
}

func (p *ArtPollReply) GetName() string {
	return string(bytes.TrimSpace(p.PortName[:]))
}

func (p *ArtPollReply) GetLongName() string {
	return string(bytes.Trim(p.LongName[:], "\u0000"))
}

func (p *ArtPollReply) GetIPAddress() net.IP {
	return net.IPv4(p.IPAddress[0], p.IPAddress[1], p.IPAddress[2], p.IPAddress[3])
}

func (p *ArtPollReply) GetSocket() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", p.GetIPAddress(), p.Port))
}

func (p *ArtPollReply) GetAvailableDMXOutputPorts() []*OutputPort {

	var list = []*OutputPort{}

	if p.Style != StNode {
		logger.Default.Debug("node style mismatch")
		return list
	}

	for i := 0; i < p.GetNumberOfPorts(); i++ {

		logger.Default.Debug("searching for output port", zap.Int("index", i))

		if p.PortTypes[i] != 0x80 { // FIXME: needs better handling
			logger.Default.Debug("port type mismatch", zap.Int("index", i), zap.Uint8("type", p.PortTypes[i]))
			continue
		}

		pa, err := NewPortAddress(p.NetSwitch, p.SubSwitch, p.SwOut[i]&0b00001111)
		if err != nil {
			logger.Default.Debug("port address", zap.Int("index", i), zap.Error(err))
			continue
		}

		socket, err := p.GetSocket()
		if err != nil {
			logger.Default.Debug("socket", zap.Int("index", i), zap.Error(err))
			continue
		}

		list = append(list, &OutputPort{
			Name:        p.GetLongName(),
			Socket:      socket,
			PortAddress: pa,
		})

	}

	return list
}

func (p *ArtPollReply) GetNumberOfPorts() int {
	return int(p.NumPortsHi)<<8 | int(p.NumPortsLo)
}

func (p *ArtPollReply) UnmarshalBinary(data []byte) error {
	if len(data) < 197 {
		return fmt.Errorf("ArtPollReply expects at least 197 bytes of data, got %d bytes", len(data))
	}
	buf := bytes.NewReader(common.PadToSize(data, p))
	return binary.Read(buf, binary.LittleEndian, p)
}

func (p *ArtPollReply) MarshalBinary() ([]byte, error) {
	buf, err := NewArtNetHeader(OpPollReply).MarshalBuffer()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.LittleEndian, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
