package eval

import (
	"io"
)

type exprKind byte

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
}

func newLiteralExpr(v token) *literalExpr {
	return &literalExpr{
		v: v,
	}
}

func (l *literalExpr) kind() exprKind {
	return Literal
}

type unaryExpr struct {
	op token
	e  expr
}

func newUnaryExpr(op token, e expr) *unaryExpr {
	return &unaryExpr{
		op: op,
		e:  e,
	}
}

func (l *unaryExpr) kind() exprKind {
	return Unary
}

type binaryExpr struct {
	l  expr
	op token
	r  expr
}

func newBinaryExpr(l expr, op token, r expr) *binaryExpr {
	return &binaryExpr{
		l:  l,
		op: op,
		r:  r,
	}
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
	tok, err := p.peek()
	if err == nil {
		return nil, newCompileError("expression expected", tok.s)
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
		left = newBinaryExpr(left, op, right)
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
		left = newBinaryExpr(left, op, right)
	}
	return left, nil
}

func parseUnary(p *parseStream) (expr, error) {
	tok, err := p.peek()
	if err != nil {
		return nil, newCompileError("expression expected", tok.s)
	}
	if tok.k != Plus && tok.k != Minus {
		return parseLit(p)
	}
	p.t.next()
	expr, err := parseUnary(p)
	if err != nil {
		return nil, err
	}
	return newUnaryExpr(tok, expr), nil
}

func parseLit(p *parseStream) (expr, error) {
	tok, err := p.next()
	if err != nil {
		return nil, err
	}
	switch tok.k {
	case Number:
		return newLiteralExpr(tok), nil
	case LeftParen:
		expr, err := parseExpression(p)
		if err != nil {
			if err == io.EOF {
				return nil, newCompileError("expression expected", tok.s)
			}
			return nil, err
		}
		_, err = p.expect(RightParen)
		if err != nil {
			return nil, err
		}
		return expr, nil
	default:
		return nil, newCompileError("expression expected", tok.s)
	}
}
