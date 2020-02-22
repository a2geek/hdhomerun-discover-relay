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
func Discover(rc *ipv4.RawConn) (net.IP, error) {
	if err := sendPacket(); err != nil {
		return nil, err
	}

	ch := waitForResponse(rc)

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case resp := <-ch:
		return resp.ip, resp.err
	case <-timer.C:
		return nil, fmt.Errorf("timed out waiting for response")
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

type responsePacket struct {
	ip  net.IP
	err error
}

func waitForResponse(rc *ipv4.RawConn) <-chan *responsePacket {
	ch := make(chan *responsePacket, 1)

	go func() {
		defer close(ch)
		data := make([]byte, 2048)
		for {
			_, ip, err := rc.ReadFromIP(data)
			if err != nil {
				ch <- &responsePacket{nil, err}
				return
			}
			udp, err := packet.BytesToUDP(data)
			if err != nil {
				ch <- &responsePacket{nil, err}
				return
			}
			fmt.Printf("packet with src=%d, dst=%d, ip=%v\n", udp.SrcPort(), udp.DstPort(), ip)
			if udp.SrcPort() != 65001 {
				continue
			}
			ch <- &responsePacket{ip.IP, nil}
			return
		}
	}()

	return ch
}
