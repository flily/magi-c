package csyntax

import (
	"testing"
)

func TestIntegerWrite(t *testing.T) {
	cases := []struct {
		value    *Integer
		expected string
	}{
		{NewIntegerLiteral(42), "42"},
		{NewHexIntegerLiteralUpper(255), "0xFF"},
		{NewHexIntegerLiteralLower(255), "0xff"},
		{NewOctalIntegerLiteral(64), "0100"},
	}

	for _, c := range cases {
		checkInterfaceCodeElement(c.value)
		checkInterfaceExpression(c.value)
		checkOutputOnStyle(t, testStyle1, c.expected, c.value)
	}
}
