package file

import (
	"bufio"
	"os"

	"github.com/russtone/utils/iter"
)

// LinesIterator represents file lines iterator.
type linesIterator struct {
	path       string
	lines      chan string
	linesCount int
	scanner    *bufio.Scanner
	file       *os.File
}

// NewLinesIterator returns new file reader.
func NewLinesIterator(path string) (iter.Iterator, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	linesCount, err := linesCount(file)
	if err != nil {
		return nil, err
	}

	if _, err = file.Seek(0, os.SEEK_SET); err != nil {
		return nil, err
	}

	return &linesIterator{
		path:       path,
		lines:      make(chan string),
		linesCount: linesCount,
		file:       file,
		scanner:    bufio.NewScanner(file),
	}, nil
}

// Next returns next read line.
func (r *linesIterator) Next(line *string) bool {

	for r.scanner.Scan() {
		*line = r.scanner.Text()
		return true
	}

	return false
}

// Close closes file reader.
func (r *linesIterator) Close() error {
	return r.file.Close()
}

// Reset resets file reader.
func (r *linesIterator) Reset() {
	r.file.Seek(0, os.SEEK_SET)
	r.scanner = bufio.NewScanner(r.file)
}

// Count returns lines count in the reader.
func (r *linesIterator) Count() uint64 {
	return uint64(r.linesCount)
}
