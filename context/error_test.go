package context

import (
	"strings"
	"testing"
)

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
