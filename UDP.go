package main

import (
	"encoding/binary"
	"fmt"
)

func BytesToUDP(data []byte) (UDP, error) {
	if len(data) < 8 {
		return nil, fmt.Errorf("Unable to create UDP packet wrapper; length must be at least 8 bytes")
	}
	return UDP(data), nil
}

// UDP is a quick and dirty wrapper around the bytes in a UDP packet.
// The assumption is the packet is valid.
type UDP []byte

func (u UDP) SrcPort() uint16 {
	return binary.BigEndian.Uint16(u[0:2])
}

func (u UDP) DstPort() uint16 {
	return binary.BigEndian.Uint16(u[2:4])
}

func (u UDP) Length() uint16 {
	return binary.BigEndian.Uint16(u[4:6])
}
