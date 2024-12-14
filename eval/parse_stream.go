package eval

import (
	"fmt"
	"io"
)

type parseStream struct {
	t *tokenStream
}

func newParseStream(r io.Reader) *parseStream {
	return &parseStream{
		t: newTokenStream(r),
	}
}

func (p *parseStream) peek() (token, error) {
	return p.t.peek()
}

func (p *parseStream) next() (token, error) {
	return p.t.next()
}

func (p *parseStream) expect(k tokenKind) (token, error) {
	tok, err := p.t.next()
	if err != nil || tok.k != k {
		msg := fmt.Sprintf("expected %s", k.string())
		return tok, newCompileError(msg, tok.s)
	}
	return tok, nil
}
