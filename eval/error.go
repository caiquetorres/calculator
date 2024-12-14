package eval

import "fmt"

type CompilerError struct {
	msg string
	s   Span
}

func newCompileError(msg string, s Span) *CompilerError {
	return &CompilerError{
		msg: msg,
		s:   s,
	}
}

func (c *CompilerError) Error() string {
	return fmt.Sprintf("%s at position %d", c.msg, c.s.s)
}

func (c *CompilerError) Span() Span {
	return c.s
}
