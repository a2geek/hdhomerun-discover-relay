package cmd

import (
	"fmt"
	"net"
	"strings"
	"time"

	"hdhomerun-discover-relay/hdhr"

	"golang.org/x/net/ipv4"
)

type DiscoverCommand struct {
}

func (cmd DiscoverCommand) Execute(args []string) error {
	pc, err := net.ListenPacket("ip4:udp", "")
	if err != nil {
		return err
	}
	defer pc.Close()
	pc.SetReadDeadline(time.Now().Add(time.Second * 1))

	rc, err := ipv4.NewRawConn(pc)
	if err != nil {
		return err
	}
	rc.SetControlMessage(ipv4.FlagDst, true)
	rc.SetControlMessage(ipv4.FlagSrc, true)

	ips, err := hdhr.Discover(rc)
	if err != nil {
		return err
	}

	switch len(ips) {
	case 0:
		fmt.Printf("No HDHomeRun(s) found!\n")
	case 1:
		fmt.Printf("%d HDHomeRun found at %s.\n", len(ips), strings.Join(ipToString(ips), ","))
	default:
		fmt.Printf("%d HDHomeRuns found at %s.\n", len(ips), strings.Join(ipToString(ips), ", "))
	}
	return nil
}

func ipToString(ips []net.IP) []string {
	s := make([]string, 0)
	for _, ip := range ips {
		s = append(s, ip.String())
	}
	return s
}
