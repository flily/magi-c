package ast

import (
	"testing"
)

func TestTokenTypeMethods(t *testing.T) {
	cases := []struct {
		t          TokenType
		isOperator bool
		s          string
	}{
		{Invalid, false, "Illegal"},
		{notUsedToken, false, "<Token 1>"},
		{True, false, "true"},
		{Plus, true, "+"},
	}

	for _, c := range cases {
		isOp := c.t.IsOperator()
		if isOp != c.isOperator {
			t.Errorf("TokenType(%d).IsOperator() == %v, expected %v", c.t, isOp, c.isOperator)
		}

		got := c.t.String()
		if got != c.s {
			t.Errorf("TokenType(%d).String() == %q, expected %q", c.t, got, c.s)
		}
	}
}

func TestGetKeywordTokenType(t *testing.T) {
	cases := []struct {
		s        string
		expected TokenType
	}{
		{"if", If},
		{"else", Else},
		{"iff", Invalid},
	}

	for _, c := range cases {
		got := GetKeywordTokenType(c.s)
		if got != c.expected {
			t.Errorf("GetKeywordTokenType(%q) == %d, expected %d", c.s, got, c.expected)
		}
	}
}
