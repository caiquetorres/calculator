package eval

import (
	"bufio"
	"io"
	"unicode"
)

type tokenStream struct {
	c   token
	e   error
	s   uint
	ptr uint
	r   *bufio.Reader
}

func newTokenStream(r io.Reader) *tokenStream {
	ts := &tokenStream{
		ptr: 0,
		s:   0,
		r:   bufio.NewReader(r),
	}
	ts.c, ts.e = ts.get()
	return ts
}

func (t *tokenStream) peekByte() (byte, error) {
	data, err := t.r.Peek(1)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

func (t *tokenStream) nextByte() (byte, error) {
	t.ptr++
	return t.r.ReadByte()
}

func (t *tokenStream) next() (token, error) {
	c, e := t.c, t.e
	t.c, t.e = t.get()
	return c, e
}

func (t *tokenStream) peek() (token, error) {
	return t.c, t.e
}

func (t *tokenStream) get() (token, error) {
	t.skipWhitespace()
	t.s = t.ptr
	ch, err := t.nextByte()
	if err != nil {
		return t.newToken(Bad), err
	}
	switch ch {
	case '+':
		return t.newToken(Plus), nil
	case '-':
		return t.newToken(Minus), nil
	case '*':
		return t.newToken(Star), nil
	case '/':
		return t.newToken(Slash), nil
	case '(':
		return t.newToken(LeftParen), nil
	case ')':
		return t.newToken(RightParen), nil
	default:
		if unicode.IsNumber(rune(ch)) {
			return t.tokNumber(ch)
		} else {
			return t.newToken(Bad), newCompileError("bad token", t.newSpan())
		}
	}
}

func (t *tokenStream) newSpan() Span {
	return Span{s: uint32(t.s), l: uint16(t.ptr - t.s)}
}

func (t *tokenStream) newToken(kind tokenKind) token {
	return token{s: t.newSpan(), k: kind}
}

func (t *tokenStream) skipWhitespace() {
	for {
		ch, err := t.peekByte()
		if err != nil {
			break
		}
		if !unicode.IsSpace(rune(ch)) {
			break
		}
		t.nextByte()
	}
}

func (t *tokenStream) tokNumber(_ byte) (token, error) {
	for {
		ch, err := t.peekByte()
		if err != nil {
			break
		}
		if !unicode.IsNumber(rune(ch)) {
			break
		}
		t.nextByte()
	}
	ch, err := t.peekByte()
	if err != nil {
		return t.newToken(Number), nil
	}
	if ch == '.' {
		t.nextByte() // '.'
		ch, err := t.peekByte()
		if err != nil || !unicode.IsNumber(rune(ch)) {
			return t.newToken(Bad), newCompileError("bad token", t.newSpan())
		}
		for {
			ch, err := t.peekByte()
			if err != nil {
				break
			}
			if !unicode.IsNumber(rune(ch)) {
				break
			}
			t.nextByte()
		}
	}
	return t.newToken(Number), nil
}
