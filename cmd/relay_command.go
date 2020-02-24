package cmd

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sync"

	"hdhomerun-discover-relay/hdhr"
	"hdhomerun-discover-relay/packet"
	"hdhomerun-discover-relay/util"

	"golang.org/x/net/ipv4"
)

type RelayCommand struct {
	Args struct {
		SourceCidr string `positional-arg-name:"cidr" description:"Source CIDR for application looking for HDHomeRun"`
	} `positional-args:"yes" required:"yes"`

	count      int
	mutex      sync.Mutex
	sourceCIDR *net.IPNet
	hdhrIPs    []net.IP
}

func (cmd *RelayCommand) Execute(args []string) error {
	var err error
	_, cmd.sourceCIDR, err = net.ParseCIDR(cmd.Args.SourceCidr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Source CIDR = %v\n", cmd.sourceCIDR)

	// listen to incoming udp packets (all interfaces)
	netpc, err := net.ListenPacket("ip4:udp", "")
	if err != nil {
		log.Fatal(err)
	}
	defer netpc.Close()

	rc, err := ipv4.NewRawConn(netpc)
	if err != nil {
		log.Fatal(err)
	}
	rc.SetControlMessage(ipv4.FlagDst, true)
	rc.SetControlMessage(ipv4.FlagSrc, true)

	cmd.hdhrIPs, err = hdhr.Discover(rc)
	if err != nil {
		log.Fatal(err)
	}
	if len(cmd.hdhrIPs) == 0 {
		log.Fatalf("No HDHomeRun discovered!")
	}
	fmt.Printf("HDHomeRun IP(s): %s\n", util.IpToString(cmd.hdhrIPs))

	fmt.Println("Starting...")
	for {
		buf := make([]byte, 10240)
		h, payload, cm, err := rc.ReadFrom(buf)
		if err != nil {
			continue
		}
		go cmd.serve(rc, h, payload, cm, buf[:h.TotalLen])
	}
}

func (cmd *RelayCommand) serve(rc *ipv4.RawConn, h *ipv4.Header, payload []byte, cm *ipv4.ControlMessage, buf []byte) {
	cmd.mutex.Lock()
	defer cmd.mutex.Unlock()

	if !cmd.sourceCIDR.Contains(h.Src) {
		fmt.Printf("Ignoring %s\n", h.Src.String())
		return
	}

	packet, err := packet.BytesToUDP(payload)
	if err != nil {
		log.Fatal(err)
	}
	if packet.DstPort() != 65001 {
		return
	}

	fmt.Printf("packet #%d, header=%v, control=%v, udp=%v\nData:\n%s",
		cmd.count, h, cm, packet, hex.Dump(buf))
	cmd.count++

	// Yeah, hack
	for _, ip := range cmd.hdhrIPs {
		copy(buf[16:], ip.To4())
		fmt.Printf("Redirecting to %v\nNew Packet:\n%s\n\n", ip, hex.Dump(buf))
		_, err = rc.WriteToIP(buf, &net.IPAddr{ip, ""})
		if err != nil {
			log.Fatal(err)
		}
	}
}
