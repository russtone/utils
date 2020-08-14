package iprange_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/russtone/utils/iprange"
)

func TestInvalid(t *testing.T) {
	tests := []struct {
		rng string
	}{
		{"300.244.42.65"},
		{"127.0.0.1/33"},
		{"10.10.10.10-1.1.1.1"},
		{"10.10.10.10-10.10.10.256"},
		{"1.1.30-1.1"},
		{"1.1.1-256.1"},
	}

	for _, tt := range tests {
		t.Run(tt.rng, func(t *testing.T) {
			_, err := iprange.Parse(tt.rng)
			require.Error(t, err)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		rng      string
		ip       string
		contains bool
	}{
		// One IP
		{"104.244.42.65", "104.244.42.65", true},
		{"87.250.250.242", "74.125.131.139", false},

		// CIDR
		{"87.240.129.133/24", "87.240.129.211", true},
		{"87.240.129.133/24", "87.240.130.211", false},
		{"140.82.118.4/22", "140.82.119.4", true},
		{"140.82.118.4/22", "140.83.119.4", false},

		// IPv4 dash range
		{"35.231.145.151-35.231.150.10", "35.231.148.1", true},
		{"35.231.145.151-35.231.150.10", "35.231.150.10", true},
		{"35.231.145.151-35.231.150.10", "35.231.145.151", true},
		{"35.231.145.151-35.231.150.10", "35.231.140.100", false},
		{"35.231.145.151-35.231.150.10", "35.231.150.100", false},

		// IPv4 octet dash/asterisk range
		{"104.16.99-100.52-55", "104.16.99.53", true},
		{"104.16.99-100.52-55", "104.16.99.52", true},
		{"104.16.99-100.52-55", "104.16.99.55", true},
		{"104.16.99-100.52-55", "104.16.100.55", true},
		{"104.16.99-100.52-55", "104.16.100.52", true},
		{"104.16.99-100.52-55", "104.16.98.2", false},
		{"104.16.99-100.52-55", "104.16.99.50", false},
		{"104.16.99.*", "104.16.99.0", true},
		{"104.16.99.*", "104.16.99.255", true},
		{"104.16.99.*", "104.16.98.50", false},
	}

	for _, tt := range tests {
		t.Run(tt.rng, func(t *testing.T) {
			r, err := iprange.Parse(tt.rng)
			require.NoError(t, err)

			contains := r.Contains(net.ParseIP(tt.ip))
			assert.Equal(t, tt.contains, contains)
		})
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		rng   string
		count uint64
	}{
		// One IP
		{"104.244.42.65", 1},
		{"87.250.250.242", 1},

		// CIDR
		{"87.240.129.133/24", 256},

		// IPv4 dash range
		{"35.231.145.151-35.231.145.200", 50},

		// IPv4 octet dash/asterisk range
		{"104.16.99-100.52-55", 8},
		{"1.1.1.*", 256},
		{"1.1.*.*", 256 * 256},
	}

	for _, tt := range tests {
		t.Run(tt.rng, func(t *testing.T) {
			r, err := iprange.Parse(tt.rng)
			require.NoError(t, err)

			assert.Equal(t, tt.count, r.Count())
		})
	}
}
