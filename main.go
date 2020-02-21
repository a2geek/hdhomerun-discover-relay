package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	"golang.org/x/net/ipv4"
)

var count = 0
var mutex = sync.Mutex{}
var sourceCIDR *net.IPNet
var sourceAddr net.Addr
var sourceIP net.IP
var sourceIfIndex int
var targetCIDR *net.IPNet
var targetAddr *net.UDPAddr
var ignoreIPs []string

func main() {
	var err error
	_, sourceCIDR, err = net.ParseCIDR("192.168.123.0/24")
	if err != nil {
		log.Fatal(err)
	}
	_, targetCIDR, err = net.ParseCIDR("192.168.5.0/24")
	if err != nil {
		log.Fatal(err)
	}
	targetAddr, err = net.ResolveUDPAddr("udp", "192.168.5.255:65001")
	if err != nil {
		log.Fatal(err)
	}
	ignoreIPs = strings.Split("192.168.5.254", ",")

	fmt.Printf("Source CIDR = %v\n", sourceCIDR)
	fmt.Printf("Target CIDR = %v\n", targetCIDR)
	fmt.Printf("Target Address = %v\n", targetAddr)
	fmt.Println("Starting...")

	// listen to incoming udp packets (all interfaces)
	netpc, err := net.ListenPacket("ip4:udp", "")
	if err != nil {
		log.Fatal(err)
	}
	defer netpc.Close()

	pc, err := ipv4.NewRawConn(netpc)
	if err != nil {
		log.Fatal(err)
	}
	pc.SetControlMessage(ipv4.FlagDst, true)
	pc.SetControlMessage(ipv4.FlagSrc, true)

	for {
		buf := make([]byte, 1024)
		h, p, cm, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		go serve(pc, h, p, cm, buf[:h.TotalLen])
	}

}

func serve(pc *ipv4.RawConn, h *ipv4.Header, p []byte, cm *ipv4.ControlMessage, buf []byte) {
	mutex.Lock()
	defer mutex.Unlock()

	if contains(ignoreIPs, h.Src.String()) {
		fmt.Printf("Ignoring %s\n\n", h.Src.String())
		return
	}

	packet, err := BytesToUDP(p)
	if err != nil {
		log.Fatal(err)
	}
	if packet.DstPort() != 65001 {
		if packet.SrcPort() == 65001 {
			fmt.Printf("RETURN packet, header=%v, control=%v, udp=%v\nPayload:\n%sData:\n%s\nUDP Payload:\n%s",
				h, cm, packet, hex.Dump(p), hex.Dump(buf), hex.Dump(packet))
		}
		return
	}

	fmt.Printf("packet #%d, header=%v, control=%v, udp=%v\nPayload:\n%sData:\n%s\nUDP Payload:\n%s",
		count, h, cm, packet, hex.Dump(p), hex.Dump(buf), hex.Dump(packet))
	count++

	bcast, err := net.ResolveIPAddr("ip", "192.168.5.117")
	if err != nil {
		log.Fatal(err)
	}
	// Yeah, hack
	copy(buf[16:], bcast.IP.To4())
	fmt.Printf("Redirecting to %v\nNew Packet:\n%s\n\n", bcast, hex.Dump(buf))
	_, err = pc.WriteToIP(buf, bcast)
	if err != nil {
		log.Fatal(err)
	}

	// host, _, err := net.SplitHostPort(addr.String())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ip := net.ParseIP(host)
	// fmt.Printf("Parsed IP = %v\n", ip)

	// if contains(ignoreIPs, host) {
	// 	fmt.Printf("Ignoring %s\n\n", host)
	// } else if sourceCIDR.Contains(ip) {
	// 	sourceAddr = addr
	// 	sourceIP = ip
	// 	sourceIfIndex = cm.IfIndex
	// 	pc.WriteTo(buf, nil, targetAddr)
	// 	fmt.Printf("Relayed to %s\n\n", targetAddr)
	// } else if targetCIDR.Contains(ip) && sourceAddr != nil {
	// 	relaycm := &ipv4.ControlMessage{
	// 		//IfIndex: sourceIfIndex,
	// 		//Src: cm.Src,
	// 		Src: net.ParseIP("192.168.5.254"),
	// 		//Dst: sourceIP,
	// 		Dst: net.ParseIP("192.168.123.11"),
	// 		TTL: 64,
	// 	}
	// 	_, port, err := net.SplitHostPort(sourceAddr.String())
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	theaddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("192.168.123.11:%s", port))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("WriteTo(buf, %v, %v)\n", relaycm, theaddr)
	// 	n, err := pc.WriteTo(buf, relaycm, theaddr)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("Relayed %d bytes to %s (control=%v)\n\n", n, theaddr, relaycm)
	// 	// pc.WriteTo(buf, nil, sourceAddr)
	// 	// fmt.Printf("Relayed to %s\n\n", sourceAddr)
	// } else {
	// 	fmt.Printf("Does not match CIDR's given; skipping...\n\n")
	// }
}

func contains(a []string, s string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}
	return false
}
