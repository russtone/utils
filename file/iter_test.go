package file_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/russtone/utils/file"
)

func TestLinesIterator(t *testing.T) {

	tests := []struct {
		file  string
		count uint64
		res   []string
	}{
		{
			"lines3",
			3,
			[]string{
				"1",
				"2",
				"3",
			},
		},
		{
			"lines3_no_newline",
			3,
			[]string{
				"1",
				"2",
				"3",
			},
		},
		{
			"lines5",
			5,
			[]string{
				"1",
				"2",
				"3",
				"4",
				"5",
			},
		},
		{
			"lines10",
			10,
			[]string{
				"1",
				"2",
				"3",
				"4",
				"5",
				"6",
				"7",
				"8",
				"9",
				"10",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			it, err := file.NewLinesIterator("test/" + tt.file)
			require.NoError(t, err)

			assert.Equal(t, tt.count, it.Count())

			var line string

			res := make([]string, 0)
			for it.Next(&line) {
				res = append(res, line)
			}

			assert.Equal(t, tt.res, res)

			it.Reset()

			res2 := make([]string, 0)
			for it.Next(&line) {
				res2 = append(res2, line)
			}

			assert.Equal(t, tt.res, res2)

			assert.NoError(t, it.Close())
		})
	}
}
