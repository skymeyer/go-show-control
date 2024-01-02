package artnet

import (
	"fmt"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"

	"go.skymyer.dev/show-control/app/logger"
	"go.skymyer.dev/show-control/common"
)

var DefaultController *Controller

func NewController(network string) (*Controller, error) {

	// Parse expected Art-Net CIDR network
	_, cidr, err := net.ParseCIDR(network)
	if err != nil {
		return nil, fmt.Errorf("invalid subnet %s: %v", network, err)
	}

	// Load all interface info
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("list interfaces: %+v", err)
	}

	// Iterate over all interface in search for the Art-Net network
	var (
		ifName  string
		ifIPNet *net.IPNet
	)
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, fmt.Errorf("list addr for interface %s: %+v", i.Name, err)
		}
		for _, addr := range addrs {
			switch a := addr.(type) {
			case *net.IPNet:
				if cidr.Contains(a.IP) && cidr.Mask.String() == a.Mask.String() {
					ifName = i.Name
					ifIPNet = a
					break
				}
			default:
				continue
			}

		}
	}

	if ifName == "" {
		return nil, fmt.Errorf("no interfaces found for Art-Net network %s", network)
	}

	ifNetwork, ifBroadcast, err := common.IPv4ToNetworkBroadcast(ifIPNet)
	if err != nil {
		return nil, err
	}

	return &Controller{
		talkToMe: NewTalkToMe().
			With(TTMReplyOnChange).
			With(TTMDiagnostics).
			With(TTMDiagnosticsUnicast).
			With(TTMDisableVLC),
		estaManCode: ESTAManUnknown,
		oemCode:     OemGlobal,
		ifName:      ifName,
		ifIPNet:     *ifIPNet,
		ifNetwork:   ifNetwork,
		ifBroadcast: ifBroadcast,
		outputPorts: make(map[PortAddress]*OutputPort),
	}, nil
}

type Controller struct {
	talkToMe    TalkToMe
	estaManCode ESTAManCode
	oemCode     OemCode
	ifName      string
	ifIPNet     net.IPNet
	ifNetwork   net.IP
	ifBroadcast net.IP

	conn       net.PacketConn
	shutdownCh chan bool

	outputPorts     map[PortAddress]*OutputPort
	outputPortsLock sync.Mutex
}

func (c *Controller) Run() error {

	logger.Default.Info("artnet controller",
		zap.String("interface", c.ifName),
		zap.String("ip", c.ifIPNet.String()),
		zap.String("network", c.ifNetwork.String()),
		zap.String("broadcast", c.ifBroadcast.String()),
	)

	conn, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", ArtNetPort))
	if err != nil {
		return err
	}
	c.conn = conn

	broadcast, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", c.ifBroadcast.String(), ArtNetPort))
	if err != nil {
		panic(err)
	}

	packet := &ArtPoll{
		Flags:        c.talkToMe,
		DiagPriority: DpLow,
	}

	packet.SetTargetPortAddress(EmptyPortAddress(), EmptyPortAddress())
	packet.SetESTAManCode(c.estaManCode)
	packet.SetOemCode(c.oemCode)

	artPoll, err := packet.MarshalBinary()
	if err != nil {
		panic(err)
	}

	var (
		tickRate = 3000 * time.Millisecond
		ticker   = time.NewTicker(tickRate)
	)

	c.shutdownCh = make(chan bool)
	recvCh := make(chan []byte)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, src, err := c.conn.ReadFrom(buf)
			if err != nil {
				logger.Default.Warn("reading packet", zap.Error(err))
				return

			}

			if src.String() == fmt.Sprintf("%s:%d", c.ifIPNet.IP.String(), ArtNetPort) {
				logger.Default.Debug("skip handling our own packet")
				continue
			}

			logger.Default.Debug("packet received from", zap.String("source", src.String()), zap.Int("bytes", n))
			recvCh <- buf[:n]
		}
	}()

	go func() {
		defer ticker.Stop()
		defer logger.Default.Debug("artnet controller terminated")

		// Fore sending ArtPoll immediately
		_, err = c.conn.WriteTo(artPoll, broadcast)
		if err != nil {
			panic(err)
		}

		for {
			select {

			case <-c.shutdownCh:
				c.conn.Close()
				c.shutdownCh = nil
				return

			case <-ticker.C:
				start := time.Now()

				_, err = c.conn.WriteTo(artPoll, broadcast)
				if err != nil {
					panic(err)
				}

				execTime := time.Since(start)
				if execTime > tickRate {
					logger.Default.Warn("artnet poll out of range",
						zap.Duration("actual", execTime),
						zap.Duration("expected", tickRate),
					)
				}

				c.gcOutputPorts()

			case data := <-recvCh:
				header, packet, err := UnmarshalBinary(data)
				if err != nil {
					logger.Default.Warn("artnet unmarshal packet", zap.Error(err))
					continue
				}
				switch v := packet.(type) {
				case *ArtPollReply:
					c.handleArtPollReply(header, v)
				default:
					logger.Default.Warn("cannot handle ArtNet packet", zap.String("opcode", header.OpCode.String()))
				}
			}
		}
	}()
	return nil
}

func (c *Controller) gcOutputPorts() {
	c.outputPortsLock.Lock()
	defer c.outputPortsLock.Unlock()

	for i, port := range c.outputPorts {
		if time.Since(port.LastSeen) > 5*time.Second {
			logger.Default.Warn("artnet node output port went away",
				zap.String("name", port.Name),
				zap.String("socket", port.Socket.String()),
				zap.String("address", port.PortAddress.String()))
			delete(c.outputPorts, i)
		}
	}
}

func (c *Controller) handleArtPollReply(h *ArtNetHeader, p *ArtPollReply) {
	logger.Default.Debug("ArtPollReply received", zap.String("node", p.GetLongName()), zap.Any("data", p))

	c.outputPortsLock.Lock()
	defer c.outputPortsLock.Unlock()

	for _, port := range p.GetAvailableDMXOutputPorts() {
		if _, found := c.outputPorts[port.PortAddress]; !found {
			logger.Default.Info("artnet output port found", zap.String("name", port.Name),
				zap.String("socket", port.Socket.String()),
				zap.String("address", port.PortAddress.String()))
			c.outputPorts[port.PortAddress] = port
		}
		c.outputPorts[port.PortAddress].LastSeen = time.Now()
	}
}

func (c *Controller) Stop() error {
	if c.shutdownCh != nil {
		c.shutdownCh <- true
	}
	return nil
}

func (c *Controller) SendDMX(pa PortAddress, data [512]byte) error {
	c.outputPortsLock.Lock()
	defer c.outputPortsLock.Unlock()

	if o, ok := c.outputPorts[pa]; ok {

		o.Sequence++

		dmx := &ArtDmx{Sequence: o.Sequence}
		dmx.SetPortAddress(pa)
		dmx.SetData(data)

		data, err := dmx.MarshalBinary()
		if err != nil {
			return err
		}

		if _, err := c.conn.WriteTo(data, o.Socket); err != nil {
			return err
		}

	}

	return nil
}

type OutputPort struct {
	Name        string
	LastSeen    time.Time
	PortAddress PortAddress
	Socket      *net.UDPAddr
	Sequence    uint8
}
