package eval

import "io"

type tokenKind byte

func (t tokenKind) string() string {
	switch t {
	case Number:
		return "number"
	case Plus:
		return "plus"
	case Minus:
		return "minus"
	case Star:
		return "star"
	case Slash:
		return "slash"
	case LeftParen:
		return "left parenthesis"
	case RightParen:
		return "right parenthesis"
	case Bad:
		return "bad"
	default:
		return "unknown"
	}
}

const (
	Number tokenKind = iota

	Plus
	Minus
	Star
	Slash

	LeftParen
	RightParen

	Bad
)

type token struct {
	k tokenKind
	s Span
}

func (t *token) textContent(r io.ReadSeeker) (string, error) {
	return t.s.textContent(r)
}
