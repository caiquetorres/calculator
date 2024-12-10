package eval

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"42", 42, false},
		{"0", 0, false},

		{"+42", 42, false},
		{"-42", -42, false},

		{"1 + 1", 2, false},
		{"10 - 5", 5, false},
		{"2 * 3", 6, false},
		{"8 / 2", 4, false},

		{"1 + 2 * 3", 7, false},
		{"(1 + 2) * 3", 9, false},
		{"10 + 5 * 2 - 8 / 4", 18, false},
		{"-3 * (2 + 4) / 3", -6, false},
		{"((1 + 2) * (3 + 4)) - 5", 16, false},

		{"", 0, true},
		{"abc", 0, true},
		{"1 / 0", 0, true},
		{"2 ** 3", 0, true},
		{"1 +", 0, true},
		{"1 + 2 3", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Eval(strings.NewReader(tc.input))
			if tc.hasError {
				assert.Error(t, err, "Expected an error for input: %s", tc.input)
			} else {
				assert.NoError(t, err, "Did not expect an error for input: %s", tc.input)
				assert.Equal(t, tc.expected, result, "Unexpected result for input: %s", tc.input)
			}
		})
	}
}
