package iprange_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/russtone/utils/iprange"
)

func TestIterator(t *testing.T) {
	tests := []struct {
		rng []string
		res []string
	}{
		// One IP
		{
			[]string{"104.244.42.65"},
			[]string{"104.244.42.65"},
		},

		// CIDR
		{
			[]string{"87.240.129.133/30"},
			[]string{
				"87.240.129.132",
				"87.240.129.133",
				"87.240.129.134",
				"87.240.129.135",
			},
		},
		{
			[]string{"87.240.129.133/32"},
			[]string{
				"87.240.129.133",
			},
		},

		// IPv4 dash range
		{
			[]string{"35.231.145.10-35.231.145.13"},
			[]string{
				"35.231.145.10",
				"35.231.145.11",
				"35.231.145.12",
				"35.231.145.13",
			},
		},

		// IPv4 octet dash/asterisk range
		{
			[]string{"104.16.99-100.52-53"},
			[]string{
				"104.16.99.52",
				"104.16.99.53",
				"104.16.100.52",
				"104.16.100.53",
			},
		},

		// Multiple
		{
			[]string{"104.244.42.65", "87.240.129.133/30"},
			[]string{
				"104.244.42.65",
				"87.240.129.132",
				"87.240.129.133",
				"87.240.129.134",
				"87.240.129.135",
			},
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.rng, ","), func(t *testing.T) {

			rr := make([]iprange.IterableRange, 0)

			for _, rng := range tt.rng {
				r, err := iprange.Parse(rng)
				require.NoError(t, err)

				rr = append(rr, r)
			}

			it := iprange.NewIterator(rr...)

			res := make([]string, 0)

			var ip string

			for it.Next(&ip) {
				res = append(res, ip)
			}

			assert.Equal(t, tt.res, res)
		})
	}
}
