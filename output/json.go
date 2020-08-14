package output

import (
	"encoding/json"
	"io"
)

const (
	rs = 0x1E
	lf = 0x0A
)

type jsonWriter struct {
	out io.Writer
}

func JSON(out io.Writer) Writer {
	return &jsonWriter{
		out: out,
	}
}

func (p *jsonWriter) Write(value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// https://tools.ietf.org/html/rfc7464
	if _, err := p.out.Write([]byte{rs}); err != nil {
		return err
	}

	if _, err := p.out.Write(b); err != nil {
		return err
	}

	if _, err := p.out.Write([]byte{lf}); err != nil {
		return err
	}

	return nil
}
