package iprange

import (
	"encoding/binary"
	"net"
)

func inc(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func incEx(ip, lower, upper net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++

		if ip[i] <= upper[i] {
			break
		} else {
			ip[i] = lower[i]
			continue
		}
	}
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}

	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
