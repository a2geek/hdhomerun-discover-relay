package cmd

import (
	"fmt"
	"net"

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

	rc, err := ipv4.NewRawConn(pc)
	if err != nil {
		return err
	}
	rc.SetControlMessage(ipv4.FlagDst, true)
	rc.SetControlMessage(ipv4.FlagSrc, true)

	ip, err := hdhr.Discover(rc)
	if err != nil {
		return err
	}

	fmt.Printf("HDHomeRun found at %v\n", ip)
	return nil
}
