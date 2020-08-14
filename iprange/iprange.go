package iprange

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"regexp"
	"strconv"
	"strings"
)

var (
	octetRegexp      = "([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])"
	octetRangeRegexp = fmt.Sprintf(`(%[1]v(-(%[1]v))?|\*)`, octetRegexp)
	cidrRegexp       = "/(3[0-2]|2[0-9]|1[0-9]|[0-9])"
	ipRegexp         = fmt.Sprintf(`(%[1]v\.){3}%[1]v`, octetRegexp)

	ipOneRegexp       = regexp.MustCompile(fmt.Sprintf(`^(%[1]v\.){3}%[1]v$`, octetRegexp))
	ipNetRegexp       = regexp.MustCompile(fmt.Sprintf(`^%v%v$`, ipRegexp, cidrRegexp))
	ipRangeRegexp     = regexp.MustCompile(fmt.Sprintf(`^(%[1]v)-(%[1]v)$`, ipRegexp))
	ipDashRangeRegexp = regexp.MustCompile(fmt.Sprintf(`^(%[1]v\.){3}%[1]v$`, octetRangeRegexp))
)

// IPRange represents common range interface.
type IPRange interface {
	// Contains checks whether the given IP is in the range.
	Contains(net.IP) bool
	Count() uint64
	next(net.IP) net.IP
}

// ipSingle represents single IP address.
// Example: "192.168.1.1"
type ipSingle struct {
	net.IP
}

var _ IPRange = ipSingle{}

// Contains checks whether the given IP is in the range.
func (r ipSingle) Contains(ip net.IP) bool {
	return r.Equal(ip)
}

// Count returns number of IP addresses in the range.
func (r ipSingle) Count() uint64 {
	return 1
}

func (r ipSingle) next(cur net.IP) net.IP {
	ip := make(net.IP, net.IPv4len)

	if cur == nil {
		copy(ip, r.IP)
		return ip
	}

	return nil
}

// ipCIDR represents CIDR.
// Example: "192.168.1.1/24"
type ipCIDR struct {
	*net.IPNet
}

var _ IPRange = ipCIDR{}

// Contains checks whether the given IP is in the range.
func (r ipCIDR) Contains(ip net.IP) bool {
	return r.IPNet.Contains(ip)
}

// Count returns number of IP addresses in the range.
func (r ipCIDR) Count() uint64 {
	ones, bits := r.Mask.Size()
	return uint64(math.Pow(2, float64(bits-ones)))
}

func (r ipCIDR) next(cur net.IP) net.IP {
	ip := make(net.IP, net.IPv4len)

	if cur == nil {
		copy(ip, r.IP)
		return ip
	}

	copy(ip, cur)

	inc(ip)

	if r.Contains(ip) {
		return ip
	}

	return nil
}

// ipMinMax represents IP range.
// Example: "192.168.1.1-192.169.1.20"
type ipMinMax struct {
	min, max net.IP
}

var _ IPRange = ipMinMax{}

// Contains checks whether the given IP is in the range.
func (r ipMinMax) Contains(ip net.IP) bool {
	cmin := bytes.Compare(ip.To4(), r.min)
	cmax := bytes.Compare(ip.To4(), r.max)
	return cmin == 0 || cmax == 0 || (cmin == 1 && cmax == -1)
}

// Count returns number of IP addresses in the range.
func (r ipMinMax) Count() uint64 {
	min := ip2int(r.min)
	max := ip2int(r.max)
	return uint64(max - min + 1)
}

func (r ipMinMax) next(cur net.IP) net.IP {
	ip := make(net.IP, net.IPv4len)

	if cur == nil {
		copy(ip, r.min)
		return ip
	}

	copy(ip, cur)
	inc(ip)

	if r.Contains(ip) {
		return ip
	}

	return nil
}

// ipLowerUpper represents IP range.
// Example: "192.168.1-3.*"
type ipLowerUpper struct {
	lower, upper net.IP
}

var _ IPRange = ipLowerUpper{}

// Contains checks whether the given IP is in the range.
func (r ipLowerUpper) Contains(ip net.IP) bool {
	for i, oct := range ip.To4() {
		if oct < r.lower.To4()[i] || oct > r.upper.To4()[i] {
			return false
		}
	}
	return true
}

// Count returns number of IP addresses in the range.
func (r ipLowerUpper) Count() uint64 {
	count := uint64(1)

	for i := len(r.lower) - 1; i >= 0; i-- {
		count *= uint64(r.upper[i]) - uint64(r.lower[i]) + 1
	}

	return count
}

func (r ipLowerUpper) next(cur net.IP) net.IP {
	ip := make(net.IP, net.IPv4len)

	if r.upper.Equal(cur) {
		return nil
	}

	if cur == nil {
		copy(ip, r.lower)
		return ip
	}

	copy(ip, cur)
	incEx(ip, r.lower, r.upper)

	if r.Contains(ip) {
		return ip
	}

	return nil
}

// Parse parses string and return corresponding IP range.
func Parse(s string) (IPRange, error) {

	switch {

	case ipOneRegexp.MatchString(s):
		ip := net.ParseIP(s)
		if ip == nil {
			return nil, fmt.Errorf("invalid ip %q", s)
		}
		return ipSingle{ip.To4()}, nil

	case ipNetRegexp.MatchString(s):
		_, ipnet, err := net.ParseCIDR(s)
		if err != nil {
			return nil, err
		}
		return ipCIDR{ipnet}, nil

	case ipRangeRegexp.MatchString(s):
		parts := strings.Split(s, "-")

		min := net.ParseIP(parts[0]).To4()
		max := net.ParseIP(parts[1]).To4()

		if ip2int(max) < ip2int(min) {
			return nil, fmt.Errorf("invalid range %q", s)
		}

		return ipMinMax{min, max}, nil

	case ipDashRangeRegexp.MatchString(s):
		lower := make(net.IP, net.IPv4len)
		upper := make(net.IP, net.IPv4len)

		parts := strings.Split(s, ".")

		for i, part := range parts {

			if part == "*" {
				lower[i] = 0
				upper[i] = 0xff
				continue
			}

			pp := strings.Split(part, "-")

			switch len(pp) {
			case 1:
				a, _ := strconv.Atoi(pp[0])
				lower[i] = byte(a)
				upper[i] = byte(a)
			case 2:
				a, _ := strconv.Atoi(pp[0])
				b, _ := strconv.Atoi(pp[1])

				if b < a {
					return nil, fmt.Errorf("invalid dash range %q", part)
				}

				lower[i] = byte(a)
				upper[i] = byte(b)
			}
		}

		return ipLowerUpper{lower, upper}, nil
	}

	return nil, fmt.Errorf("invalid range %q", s)
}
