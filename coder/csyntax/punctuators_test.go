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
	builder, writer := makeTestWriter(testStyle1)

	err := writer.Write(0, OperatorAdd, OperatorAssign)
	if err != nil {
		t.Fatalf("WriteItems Punctuators failed: %s", err)
	}

	expected := "+="
	result := builder.String()
	if result != expected {
		t.Fatalf("WriteItems Punctuators result wrong, expected '%s', got '%s'", expected, result)
	}
}
