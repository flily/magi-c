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
		checkInterfaceCodeElement(c.op)

		result := c.op.String()
		if result != c.expected {
			t.Fatalf("Punctuator String result wrong for %d, expect '%s', got '%s'", c.op, c.expected, result)
		}
	}
}

func TestWritePunctuators(t *testing.T) {
	expected := "+="
	checkOutputOnStyle(t, testStyle1, expected, OperatorAdd, OperatorAssign)
}
