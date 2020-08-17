package dns

import (
	"fmt"
	"net"
	"strings"
)

// Subdomains returns list of subdomains of base domain.
func Subdomains(sub, base string) []string {
	domains := make([]string, 0)

	ls := strings.Count(sub, ".")
	lb := strings.Count(base, ".")

	if ls <= lb {
		return nil
	}

	parts := strings.Split(sub, ".")

	for i := 0; i < ls-lb; i++ {
		d := strings.Join(parts[i:], ".")
		domains = append(domains, d)
	}

	return domains
}

const hexDigit = "0123456789abcdef"

// PTR returns domain name for IP.
func PTR(ip net.IP) string {

	if ip.To4() != nil {
		parts := strings.Split(ip.String(), ".")
		return fmt.Sprintf("%s.%s.%s.%s.in-addr.arpa", parts[3], parts[2], parts[1], parts[0])
	}

	// Must be IPv6
	buf := make([]byte, 0, len(ip)*4+len("ip6.arpa"))

	// Add it, in reverse, to the buffer
	for i := len(ip) - 1; i >= 0; i-- {
		v := ip[i]
		buf = append(buf, hexDigit[v&0xF],
			'.',
			hexDigit[v>>4],
			'.')
	}

	buf = append(buf, "ip6.arpa"...)

	return string(buf)
}
