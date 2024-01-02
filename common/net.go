package common

import (
	"encoding/binary"
	"fmt"
	"net"
)

func IPv4ToNetworkBroadcast(n *net.IPNet) (net.IP, net.IP, error) {
	if n.IP.To4() == nil {
		return net.IP{}, net.IP{}, fmt.Errorf("input no IPv4 address")
	}
	var (
		network   = make(net.IP, len(n.IP.To4()))
		broadcast = make(net.IP, len(n.IP.To4()))
	)
	binary.BigEndian.PutUint32(network, binary.BigEndian.Uint32(n.IP.To4())&binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
	binary.BigEndian.PutUint32(broadcast, binary.BigEndian.Uint32(n.IP.To4())|^binary.BigEndian.Uint32(net.IP(n.Mask).To4()))
	return network, broadcast, nil
}
