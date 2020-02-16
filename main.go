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

func main() {
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

	fmt.Printf("packet #%d, network=%s, addr=%s\nData:\n%s\n", count, addr.Network(), addr.String(), hex.Dump(buf))
	count++
}
