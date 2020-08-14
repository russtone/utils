package term

import (
	"io"

	"golang.org/x/sys/unix"
)

// reader is an io.Reader that reads from a specific file descriptor.
type reader int

func (r reader) Read(buf []byte) (int, error) {
	return unix.Read(int(r), buf)
}

// WaitForEnter reads input from a terminal without local echo until "\n" is found.
func WaitForEnter(fd int) error {
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		return err
	}

	newState := *termios
	newState.Lflag &^= (unix.ECHO | unix.ICANON)
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &newState); err != nil {
		return err
	}

	defer unix.IoctlSetTermios(fd, ioctlWriteTermios, termios)

	return readLine(reader(fd))
}

func readLine(reader io.Reader) error {
	var buf [1]byte

	for {
		n, err := reader.Read(buf[:])
		if n > 0 {
			if buf[0] == '\n' {
				return nil
			}
			continue
		}

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}
	}
}
