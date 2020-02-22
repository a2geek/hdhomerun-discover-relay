package hdhr

import (
	"fmt"
	"hdhomerun-discover-relay/packet"
	"net"
	"time"

	"golang.org/x/net/ipv4"
)

var discoveryPacket = []byte{
	0x00, 0x02, // discover request
	0x00, 0x0c, // length
	0x01,                   // device type
	0x04,                   // length of 4
	0xff, 0xff, 0xff, 0xff, // device type wildcard
	0x02,                   // device id
	0x04,                   // length of 4
	0xff, 0xff, 0xff, 0xff, // target id wildcard
	0x73, 0xcc, 0x7d, 0x8f, // crc
}

func Discover(rc *ipv4.RawConn) (net.IP, error) {
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 65001,
	})
	defer socket.Close()

	_, err = socket.Write(discoveryPacket)
	if err != nil {
		return nil, err
	}

	ipch := make(chan net.IP, 1)
	defer close(ipch)

	errch := make(chan error, 1)
	defer close(errch)

	go func() {
		data := make([]byte, 2048)
		for {
			_, ip, err := rc.ReadFromIP(data)
			if err != nil {
				errch <- err
				return
			}
			udp, err := packet.BytesToUDP(data)
			if err != nil {
				errch <- err
				return
			}
			fmt.Printf("packet with src=%d, dst=%d, ip=%v\n", udp.SrcPort(), udp.DstPort(), ip)
			if udp.SrcPort() != 65001 {
				continue
			}
			ipch <- ip.IP
			return
		}
	}()

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case ip := <-ipch:
		return ip, nil
	case err = <-errch:
		return nil, err
	case <-timer.C:
		return nil, fmt.Errorf("timed out waiting for response")
	}
}
