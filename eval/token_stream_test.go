package eval

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenStream_Next(t *testing.T) {
	input := "+-*/()123"
	reader := strings.NewReader(input)
	ts := newTokenStream(reader)

	tests := []struct {
		expectedToken tokenKind
		expectedErr   error
	}{
		{Plus, nil},
		{Minus, nil},
		{Star, nil},
		{Slash, nil},
		{LeftParen, nil},
		{RightParen, nil},
		{Number, nil}, // '1'
		{0, io.EOF},   // End of input
	}

	for _, tt := range tests {
		tok, err := ts.next()
		assert.Equal(t, tt.expectedToken, tok.k, "Token kind mismatch")
		assert.True(t, errors.Is(err, tt.expectedErr), "Error mismatch")
	}
}

func TestTokenStream_Peek(t *testing.T) {
	input := "+-"
	reader := strings.NewReader(input)
	ts := newTokenStream(reader)

	tok, err := ts.peek()
	assert.Equal(t, Plus, tok.k, "Peek should return Plus as the first token")
	assert.NoError(t, err, "Peek should not return an error")

	tok, err = ts.next()
	assert.Equal(t, Plus, tok.k, "Next should return Plus after Peek")
	assert.NoError(t, err, "Next should not return an error after Peek")

	tok, err = ts.peek()
	assert.Equal(t, Minus, tok.k, "Peek should return Minus as the second token")
	assert.NoError(t, err, "Peek should not return an error")
}

func TestTokenStream_EOF(t *testing.T) {
	input := ""
	reader := strings.NewReader(input)
	ts := newTokenStream(reader)

	tok, err := ts.next()
	assert.Equal(t, tokenKind(0), tok.k, "Token should be empty at EOF")
	assert.ErrorIs(t, err, io.EOF, "Error should be EOF for empty input")
}

func TestTokenStream_InvalidCharacter(t *testing.T) {
	input := "?"
	reader := strings.NewReader(input)
	ts := newTokenStream(reader)

	tok, err := ts.next()
	assert.Equal(t, tokenKind(0), tok.k, "Token should be empty for invalid character")
	assert.Error(t, err, "Error should be returned for invalid character")
}
