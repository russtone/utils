package rnd_test

import (
	"testing"

	"github.com/russtone/utils/rnd"
	"github.com/stretchr/testify/assert"
)

func Test_String(t *testing.T) {
	for i := 1; i < 100; i += 10 {
		assert.Len(t, rnd.String(i), i)
	}
}
