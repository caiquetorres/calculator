package eval

import (
	"io"
)

type exprKind byte

func (e exprKind) string() string {
	switch e {
	case Literal:
		return "literal"
	case Unary:
		return "unary"
	case Binary:
		return "binary"
	default:
		return "unkdown"
	}
}

const (
	Literal exprKind = iota
	Unary
	Binary
)

type expr interface {
	kind() exprKind
}

type literalExpr struct {
	v token

	k exprKind
}

func (l *literalExpr) kind() exprKind {
	return Literal
}

type unaryExpr struct {
	op token
	e  expr

	k exprKind
}

func (l *unaryExpr) kind() exprKind {
	return Unary
}

type binaryExpr struct {
	l  expr
	op token
	r  expr

	k exprKind
}

func (l *binaryExpr) kind() exprKind {
	return Binary
}

func parse(r io.Reader) (expr, error) {
	p := newParseStream(r)
	e, err := parseExpression(p)
	if err != nil {
		return nil, err
	}
	_, err = p.peek()
	if err == nil {
		return nil, errUnexpectedTok
	}
	return e, nil
}

func parseExpression(p *parseStream) (expr, error) {
	e, err := parseTerm(p)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func parseTerm(p *parseStream) (expr, error) {
	left, err := parseFactor(p)
	if err != nil {
		return nil, err
	}
	for {
		tok, err := p.peek()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if tok.k != Plus && tok.k != Minus {
			break
		}
		op, _ := p.next()
		right, err := parseFactor(p)
		if err != nil {
			return nil, err
		}
		left = &binaryExpr{
			l:  left,
			op: op,
			r:  right,
		}
	}
	return left, nil
}

func parseFactor(p *parseStream) (expr, error) {
	left, err := parseUnary(p)
	if err != nil {
		return nil, err
	}
	for {
		tok, err := p.peek()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if tok.k != Star && tok.k != Slash {
			break
		}
		op, _ := p.next()
		right, err := parseUnary(p)
		if err != nil {
			return nil, err
		}
		left = &binaryExpr{
			l:  left,
			op: op,
			r:  right,
		}
	}
	return left, nil
}

func parseUnary(p *parseStream) (expr, error) {
	tok, err := p.peek()
	if err != nil {
		return nil, errExpressionExpected
	}
	if tok.k != Plus && tok.k != Minus {
		return parseLit(p)
	}
	p.t.next()
	expr, err := parseUnary(p)
	if err != nil {
		return nil, err
	}
	return &unaryExpr{
		op: tok,
		e:  expr,
	}, nil
}

func parseLit(p *parseStream) (expr, error) {
	tok, err := p.next()
	if err != nil {
		return nil, err
	}
	switch tok.k {
	case Number:
		return &literalExpr{
			v: tok,
			k: Literal,
		}, nil
	case LeftParen:
		expr, err := parseExpression(p)
		if err != nil {
			if err == io.EOF {
				return nil, errExpressionExpected
			}
			return nil, err
		}
		_, err = p.expect(RightParen)
		if err != nil {
			return nil, err
		}
		return expr, nil
	default:
		return nil, errExpressionExpected
	}
}
