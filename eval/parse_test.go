package eval

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLiteral(t *testing.T) {
	input := "42"
	expr, err := parse(strings.NewReader(input))

	assert.NoError(t, err, "Parsing a literal should not produce an error")
	assert.IsType(t, &literalExpr{}, expr, "Parsed expression should be a literal")
	lit := expr.(*literalExpr)
	assert.Equal(t, Literal, lit.kind(), "Literal kind should match")
	assert.Equal(t, Number, lit.v.k, "Literal token kind should be Number")
}

func TestParseUnary(t *testing.T) {
	input := "-42"
	expr, err := parse(strings.NewReader(input))

	assert.NoError(t, err, "Parsing a unary expression should not produce an error")
	assert.IsType(t, &unaryExpr{}, expr, "Parsed expression should be unary")
	unary := expr.(*unaryExpr)
	assert.Equal(t, Unary, unary.kind(), "Unary kind should match")
	assert.Equal(t, Minus, unary.op.k, "Unary operator should be Minus")
	assert.IsType(t, &literalExpr{}, unary.e, "Unary operand should be a literal")
}

func TestParseBinary(t *testing.T) {
	input := "6 * 7"
	expr, err := parse(strings.NewReader(input))

	assert.NoError(t, err, "Parsing a binary expression should not produce an error")
	assert.IsType(t, &binaryExpr{}, expr, "Parsed expression should be binary")
	binary := expr.(*binaryExpr)
	assert.Equal(t, Binary, binary.kind(), "Binary kind should match")
	assert.Equal(t, Star, binary.op.k, "Binary operator should be Star")
	assert.IsType(t, &literalExpr{}, binary.l, "Left operand should be a literal")
	assert.IsType(t, &literalExpr{}, binary.r, "Right operand should be a literal")
}

func TestParseParentheses(t *testing.T) {
	input := "(1 + 2) * 3"
	expr, err := parse(strings.NewReader(input))

	assert.NoError(t, err, "Parsing a parenthesized expression should not produce an error")
	assert.IsType(t, &binaryExpr{}, expr, "Parsed expression should be binary")

	binary := expr.(*binaryExpr)
	assert.Equal(t, Binary, binary.kind(), "Binary kind should match")
	assert.Equal(t, Star, binary.op.k, "Binary operator should be Star")

	assert.IsType(t, &binaryExpr{}, binary.l, "Left operand should be a binary expression")
	left := binary.l.(*binaryExpr)
	assert.Equal(t, Plus, left.op.k, "Inner binary operator should be Plus")
	assert.IsType(t, &literalExpr{}, left.l, "Left operand of inner binary should be a literal")
	assert.IsType(t, &literalExpr{}, left.r, "Right operand of inner binary should be a literal")
}

func TestParseNestedBinary(t *testing.T) {
	input := "1 + 2 * 3"
	expr, err := parse(strings.NewReader(input))

	assert.NoError(t, err, "Parsing a nested binary expression should not produce an error")
	assert.IsType(t, &binaryExpr{}, expr, "Parsed expression should be binary")

	binary := expr.(*binaryExpr)
	assert.Equal(t, Binary, binary.kind(), "Binary kind should match")
	assert.Equal(t, Plus, binary.op.k, "Outer binary operator should be Plus")

	assert.IsType(t, &literalExpr{}, binary.l, "Left operand should be a literal")
	assert.IsType(t, &binaryExpr{}, binary.r, "Right operand should be a binary expression")

	right := binary.r.(*binaryExpr)
	assert.Equal(t, Star, right.op.k, "Inner binary operator should be Star")
	assert.IsType(t, &literalExpr{}, right.l, "Left operand of inner binary should be a literal")
	assert.IsType(t, &literalExpr{}, right.r, "Right operand of inner binary should be a literal")
}

func TestParseInvalid(t *testing.T) {
	input := "? + 2"
	_, err := parse(strings.NewReader(input))

	assert.Error(t, err, "Parsing invalid input should produce an error")
}
