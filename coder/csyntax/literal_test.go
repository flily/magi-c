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
		{NewHexIntegerLiteral(255), "0xff"},
		{NewOctalIntegerLiteral(64), "0100"},
	}

	for _, c := range cases {
		builder, writer := makeTestWriter(KRStyle)
		err := c.value.Write(writer, 0)
		if err != nil {
			t.Fatalf("Integer Write failed, '%d': %s", c.value, err)
		}

		result := builder.String()
		if result != c.expected {
			t.Fatalf("Integer Write result wrong for value %d, expected: %s, got: %s",
				c.value, c.expected, result)
		}
	}
}
