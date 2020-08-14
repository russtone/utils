package file

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

const bufSize = 32 * 1024

// LinesCount returns count of lines in file.
func LinesCount(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return linesCount(file)
}

func linesCount(r io.Reader) (int, error) {
	buf := make([]byte, bufSize)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		// If there is no "\n" on the last line.
		if c != 0 && c < bufSize && buf[c-1] != '\n' {
			count++
		}

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// FirstLine returns first line of file.
func FirstLine(path string) (string, error) {

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			return line, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", err
}
