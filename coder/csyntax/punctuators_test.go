package csyntax

import (
	"testing"
)

func TestPunctuatorString(t *testing.T) {
	cases := []struct {
		op       Punctuator
		expected string
	}{
		{OperatorAssign, "="},
		{OperatorAdd, "+"},
		{OperatorArrow, "->"},
		{PunctuatorSemicolon, ";"},
		{0, "INVALID"},
	}

	for _, c := range cases {
		result := c.op.String()
		if result != c.expected {
			t.Fatalf("Punctuator String result wrong for %d, expect '%s', got '%s'", c.op, c.expected, result)
		}
	}
}
