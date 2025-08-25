package token

import (
	"testing"
)

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
