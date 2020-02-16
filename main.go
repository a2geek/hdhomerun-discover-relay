package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sync"
)

var count = 0
var mutex = sync.Mutex{}
var sourceCIDR *net.IPNet
var sourceAddr net.Addr
var targetCIDR *net.IPNet
var targetAddr *net.UDPAddr

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

	fmt.Printf("Source CIDR = %v\n", sourceCIDR)
	fmt.Printf("Target CIDR = %v\n", targetCIDR)
	fmt.Printf("Target Address = %v\n", targetAddr)
	fmt.Println("Starting...")

	// listen to incoming udp packets
	pc, err := net.ListenPacket("udp", ":65001")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		go serve(pc, addr, buf[:n])
	}

}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Printf("packet #%d, network=%s, addr=%s\nData:\n%s", count, addr.Network(), addr.String(), hex.Dump(buf))
	count++

	host, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		log.Fatal(err)
	}
	ip := net.ParseIP(host)
	fmt.Printf("Parsed IP = %v\n", ip)

	if sourceCIDR.Contains(ip) {
		sourceAddr = addr
		pc.WriteTo(buf, targetAddr)
		fmt.Printf("Relayed to %s\n\n", targetAddr)
	} else if targetCIDR.Contains(ip) {
		pc.WriteTo(buf, sourceAddr)
		fmt.Printf("Relayed to %s\n\n", sourceAddr)
	} else {
		fmt.Printf("Does not match CIDR's given; skipping...\n\n")
	}
}
