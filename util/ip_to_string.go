package util

import "net"

// Convert a slice of net.IP to string
func IpToString(ips []net.IP) []string {
	s := make([]string, 0)
	for _, ip := range ips {
		s = append(s, ip.String())
	}
	return s
}
