package eval

import "io"

type span struct {
	s uint32
	l uint16
}

func (s *span) textContent(r io.ReadSeeker) (string, error) {
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
