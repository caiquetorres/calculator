package eval

import (
	"fmt"
	"io"
	"strconv"
)

func Eval(r io.ReadSeeker) (float64, error) {
	expr, err := parse(r)
	if err != nil {
		return 0, err
	}
	res, err := evalExpr(r, expr)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func evalExpr(r io.ReadSeeker, e expr) (float64, error) {
	switch e.kind() {
	case Literal:
		e := e.(*literalExpr)
		t, err := e.v.textContent(r)
		if err != nil {
			return 0, nil
		}
		i, err := strconv.Atoi(t)
		if err != nil {
			return 0, nil
		}
		return float64(i), nil
	case Unary:
		e := e.(*unaryExpr)
		u, err := evalExpr(r, e.e)
		if err != nil {
			return 0, nil
		}
		opText, err := e.op.textContent(r)
		if err != nil {
			return 0, nil
		}
		switch opText {
		case "+":
			return u, nil
		case "-":
			return -u, nil
		}
	case Binary:
		e := e.(*binaryExpr)
		left, err := evalExpr(r, e.l)
		if err != nil {
			return 0, nil
		}
		right, err := evalExpr(r, e.r)
		if err != nil {
			return 0, nil
		}
		opText, err := e.op.textContent(r)
		if err != nil {
			return 0, nil
		}
		switch opText {
		case "+":
			return left + right, nil
		case "-":
			return left - right, nil
		case "*":
			return left * right, nil
		case "/":
			if right == 0 {
				return 0, fmt.Errorf("invalid division bg=y 0")
			}
			return left / right, nil
		}
	}
	return 0, nil
}
