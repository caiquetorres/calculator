package eval

import (
	"io"
)

type Span struct {
	s uint32
	l uint16
}

func (s *Span) Start() uint32 {
	return s.s
}

func (s *Span) textContent(r io.ReadSeeker) (string, error) {
	buf := make([]byte, s.l)
	_, err := r.Seek(int64(s.s), io.SeekStart)
	if err != nil {
		return "", err
	}
	_, err = r.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
