package artnet

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"go.skymyer.dev/show-control/app/logger"
	"go.skymyer.dev/show-control/dmx"
	"go.skymyer.dev/show-control/dmx/driver"
)

func init() {
	driver.Register(DRIVER_NAME, NewArtNetNode)
}

const (
	DRIVER_NAME     = "artnet"
	REFRESH_RATE_MS = 23 // Note: 44Hz or 22,72 ms, fetch from poll if available from node
)

func NewArtNetNode(device string) (driver.Driver, error) {
	n := &ArtNetNode{
		device: device,
	}
	if err := n.Open(); err != nil {
		return nil, err
	}
	return n, nil
}

type ArtNetNode struct {
	device      string
	portAddress PortAddress
	output      *dmx.Frame
	shutdownCh  chan bool
}

func (n *ArtNetNode) SetUniverse(universe int, output *dmx.Frame) error {
	pa, err := ParsePortAddress(fmt.Sprintf("%s-%d", n.device, universe))
	if err != nil {
		return err
	}
	n.portAddress = pa
	n.output = output
	return nil
}

func (n *ArtNetNode) Open() error {
	// TODO - add discovery
	return nil
}

func (n *ArtNetNode) Close() error {
	return nil
}

func (n *ArtNetNode) Run() error {
	tickRate := REFRESH_RATE_MS * time.Millisecond
	ticker := time.NewTicker(tickRate)
	n.shutdownCh = make(chan bool)
	go func() {
		defer ticker.Stop()
		defer logger.Default.Debug("artnet node terminated")
		for {
			select {
			case <-n.shutdownCh:
				n.shutdownCh = nil
				return
			case <-ticker.C:
				start := time.Now()
				n.send(n.output)
				execTime := time.Since(start)
				if execTime > tickRate {
					logger.Default.Warn("artnet node apply exec out of range",
						zap.Duration("actual", execTime),
						zap.Duration("expected", tickRate),
					)
				}
			}
		}
	}()
	return nil
}

func (n *ArtNetNode) Stop() error {
	if n.shutdownCh != nil {
		n.shutdownCh <- true
	}

	// clear out all outputs
	return n.send(dmx.NewDMX512Frame())
}

func (n *ArtNetNode) send(frame *dmx.Frame) error {
	if DefaultController != nil {
		var data [512]byte
		copy(data[:], frame.ReadSlots())
		return DefaultController.SendDMX(n.portAddress, data)
	} else {
		logger.Default.Warn("artnet controller not initialized yet")
	}
	return nil
}
