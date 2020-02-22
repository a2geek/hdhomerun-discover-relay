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

// Discover sends the HDHomeRun discovery packet and waits for a response.
func Discover(rc *ipv4.RawConn) ([]net.IP, error) {
	if err := sendPacket(); err != nil {
		return nil, err
	}

	active := true
	ipch, errch := waitForResponse(rc, &active)

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	ips := make([]net.IP, 0)

	// collect IP addresses for timer duration
	for {
		select {
		case ip := <-ipch:
			ips = append(ips, ip)
		case err := <-errch:
			return nil, err
		case <-timer.C:
			active = false
			if len(ips) == 0 {
				return nil, fmt.Errorf("timed out waiting for response")
			} else {
				return ips, nil
			}
		}
	}
}

func sendPacket() error {
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 65001,
	})
	defer socket.Close()

	_, err = socket.Write(discoveryPacket)
	return err
}

func waitForResponse(rc *ipv4.RawConn, active *bool) (<-chan net.IP, <-chan error) {
	ipch := make(chan net.IP, 1)
	errch := make(chan error, 1)

	go func() {
		defer close(ipch)
		defer close(errch)

		data := make([]byte, 2048)
		for *active {
			_, ip, err := rc.ReadFromIP(data)
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				continue
			} else if err != nil {
				errch <- err
				return
			}
			udp, err := packet.BytesToUDP(data)
			if err != nil {
				errch <- err
				return
			}
			if udp.SrcPort() != 65001 {
				continue
			}
			ipch <- ip.IP
		}
	}()

	return ipch, errch
}
