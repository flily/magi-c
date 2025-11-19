package tokenizer

import (
	"testing"
)

func TestValidIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"_", true},
		{"valid_identifier", true},
		{"validIdentifier123", true},
		{"ValidIdentifier", true},
		{"_underscore", true},
		{"123invalid", false},
		{"invalid-char!", false},
	}

	for _, test := range tests {
		result := IsValidIdentifier(test.input)
		if result != test.expected {
			t.Errorf("IsValidIdentifier(%q) = %v; want %v", test.input, result, test.expected)
		}
	}
}
