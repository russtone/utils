package slice

import (
	"net"
	"sort"
)

// ContainsString returns true if string in slice and false otherwise.
func ContainsString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// ContainsIP returns true if ip in slice and false otherwise.
func ContainsIP(ips []net.IP, ip net.IP) bool {
	for _, i := range ips {
		if i.String() == ip.String() {
			return true
		}
	}
	return false
}

// Dedup returns slice of strings without duplicates.
func DedupStrings(items []string) []string {
	if len(items) == 0 {
		return []string{}
	}

	sort.Strings(items)

	j := 0
	for i := 1; i < len(items); i++ {
		if items[j] == items[i] {
			continue
		}

		j++

		items[j] = items[i]
	}

	return items[:j+1]
}
