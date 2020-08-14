package iter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/russtone/utils/iter"
)

func Test_Slice(t *testing.T) {
	tests := []struct {
		slice []string
	}{
		{[]string{"one", "two", "three"}},
	}

	for _, vec := range tests {
		it := iter.Slice(vec.slice)

		assert.Equal(t, uint64(len(vec.slice)), it.Count())

		var s string
		res := make([]string, 0)

		for it.Next(&s) {
			res = append(res, s)
		}

		assert.Equal(t, vec.slice, res)

		assert.NoError(t, it.Close())
	}
}
