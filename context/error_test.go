package context

import (
	"strings"
	"testing"
)

func TestErrorLevelNames(t *testing.T) {
	tests := []struct {
		level    ErrorLevel
		expected string
	}{
		{Note, "note"},
		{Warning, "warning"},
		{Error, "error"},
	}

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Fatalf("error level name mismatch for level %d: expected %s, got %s",
				test.level, test.expected, test.level.String())
		}
	}
}

func TestErrorMessage(t *testing.T) {
	fd := createTestFile1()

	line := fd.LineContext(3)
	ctx := line.Mark(7, 14)

	err := ctx.Error("the quick brown fox")
	expected := strings.Join([]string{
		"   4:   sed do eiusmod tempor incididunt",
		"               ^^^^^^^",
		"               the quick brown fox",
	}, "\n")

	if err.Error() != expected {
		t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, err.Error())
	}
}
