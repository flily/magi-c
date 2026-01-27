package context

import (
	"testing"

	"strings"
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

	err := ctx.Error("the quick brown fox").With("jumps over the lazy dog")
	expected := strings.Join([]string{
		"example.txt:4:8: error: the quick brown fox",
		"    4 | sed do eiusmod tempor incididunt",
		"      |        ^^^^^^^",
		"      |        jumps over the lazy dog",
	}, "\n")

	if err.Error() != expected {
		t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, err.Error())
	}
}

func TestDiagnosticLevelString(t *testing.T) {
	fd := createTestFile1()

	line := fd.LineContext(3)
	ctx := line.Mark(7, 14)

	{
		err := ctx.Note("the quick brown fox").With("jumps over the lazy dog")
		expected := strings.Join([]string{
			"example.txt:4:8: note: the quick brown fox",
			"    4 | sed do eiusmod tempor incididunt",
			"      |        ^^^^^^^",
			"      |        jumps over the lazy dog",
		}, "\n")

		if err.Error() != expected {
			t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, err.Error())
		}
	}

	{
		err := ctx.Warning("the quick brown fox").With("jumps over the lazy dog")
		expected := strings.Join([]string{
			"example.txt:4:8: warning: the quick brown fox",
			"    4 | sed do eiusmod tempor incididunt",
			"      |        ^^^^^^^",
			"      |        jumps over the lazy dog",
		}, "\n")

		if err.Error() != expected {
			t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, err.Error())
		}
	}

	{
		err := ctx.Error("the quick brown fox").With("jumps over the lazy dog")
		expected := strings.Join([]string{
			"example.txt:4:8: error: the quick brown fox",
			"    4 | sed do eiusmod tempor incididunt",
			"      |        ^^^^^^^",
			"      |        jumps over the lazy dog",
		}, "\n")

		if err.Error() != expected {
			t.Fatalf("error message mismatch, expected:\n%s\ngot:\n%s", expected, err.Error())
		}
	}
}

func TestDiagnosticComboMessage(t *testing.T) {
	fd := createTestFile1()

	line1 := fd.LineContext(3)
	ctx1 := line1.Mark(7, 14)

	line2 := fd.LineContext(5)
	ctx2 := line2.Mark(11, 16)

	err := ctx1.Error("the quick brown fox").
		With("jumps over the lazy dog").
		For(ctx2.Note("lorem ipsum").With("dolor sit amet"))

	expected := strings.Join([]string{
		"example.txt:4:8: error: the quick brown fox",
		"    4 | sed do eiusmod tempor incididunt",
		"      |        ^^^^^^^",
		"      |        jumps over the lazy dog",
		"example.txt:6:12: note: lorem ipsum",
		"    6 | ut enim ad minim veniam",
		"      |            ^^^^^",
		"      |            dolor sit amet",
	}, "\n")

	if err.Error() != expected {
		t.Fatalf("diagnostic combo message mismatch, expected:\n%s\ngot:\n%s", expected, err.Error())
	}
}
