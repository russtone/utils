package file_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/russtone/utils/file"
)

func TestLinesCount(t *testing.T) {
	tests := []struct {
		file  string
		count int
	}{
		{"lines3_no_newline", 3},
		{"lines3", 3},
		{"lines5", 5},
		{"lines10", 10},
		{"lines1000", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			count, err := file.LinesCount("test/" + tt.file)
			require.NoError(t, err)
			assert.Equal(t, tt.count, count)
		})
	}
}

func TestLinesCount_Error(t *testing.T) {
	_, err := file.LinesCount("test/invalid")
	assert.Error(t, err)
}

func TestFirstLine(t *testing.T) {
	tests := []struct {
		file string
		line string
	}{
		{"lines3_no_newline", "1"},
		{"lines3", "1"},
		{"lines5", "1"},
		{"lines10", "1"},
		{"lines1000", "1"},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			line, err := file.FirstLine("test/" + tt.file)
			require.NoError(t, err)
			assert.Equal(t, tt.line, line)
		})
	}
}
