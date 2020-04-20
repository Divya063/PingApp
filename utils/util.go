package utils

import "net"

func IsIPv4(ip net.IP) bool {
	return len(ip.To4()) == net.IPv4len
}

func IsIPv6(ip net.IP) bool {
	return len(ip) == net.IPv6len
}
