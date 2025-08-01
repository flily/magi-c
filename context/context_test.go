package context

import (
	"testing"

	"bytes"
	"strings"
)

func createTestFile1() *FileContext {
	parts := [][]byte{
		[]byte("lorem ipsum dolor sit amet\n"),
		[]byte("consectetur adipiscing elit\n"),
		[]byte("\n"),
		[]byte("sed do eiusmod tempor incididunt\n"),
		[]byte("ut labore et dolore magna aliqua\n"),
		[]byte("ut enim ad minim veniam\r\n"),
		[]byte("\r\n"),
		[]byte("quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat\n"),
		[]byte("duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur\n"),
		[]byte("excepteur sint occaecat cupid atat non proident\n"),
	}

	content := bytes.Join(parts, []byte{})
	file := ReadFileData("example.txt", content)
	return file
}

func TestContextHighlightText(t *testing.T) {
	fd := createTestFile1()

	line := fd.LineContext(3)
	// 0         1         2         3         4
	// 0    5    0    5    0    5    0    5    0
	// sed do eiusmod tempor incididunt
	ctx := line.Mark(7, 14)
	ctx.Load(2, 2)

	got := ctx.HighlightText("the quick brown fox")
	expected := strings.Join([]string{
		"   2:   consectetur adipiscing elit",
		"   3:   ",
		"   4:   sed do eiusmod tempor incididunt",
		"               ^^^^^^^",
		"               the quick brown fox",
		"   5:   ut labore et dolore magna aliqua",
		"   6:   ut enim ad minim veniam",
	}, "\n")

	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestContextHighlightTextMultipleParts1(t *testing.T) {
	fd := createTestFile1()

	line1 := fd.LineContext(3)
	// 0         1         2         3         4
	// 0    5    0    5    0    5    0    5    0
	// sed do eiusmod tempor incididunt
	ctx1 := line1.Mark(7, 14)

	line2 := fd.LineContext(3)
	ctx2 := line2.Mark(15, 21)

	ctx := ctx1.Join(ctx2)
	ctx.Load(2, 2)

	got := ctx.HighlightText("the quick brown fox")
	expected := strings.Join([]string{
		"   2:   consectetur adipiscing elit",
		"   3:   ",
		"   4:   sed do eiusmod tempor incididunt",
		"               ^^^^^^^ ^^^^^^",
		"               the quick brown fox",
		"   5:   ut labore et dolore magna aliqua",
		"   6:   ut enim ad minim veniam",
	}, "\n")

	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestContextHighlightTextMultipleParts2(t *testing.T) {
	fd := createTestFile1()

	line1 := fd.LineContext(3)
	// 0         1         2         3         4
	// 0    5    0    5    0    5    0    5    0
	// sed do eiusmod tempor incididunt
	ctx1 := line1.Mark(7, 14)

	line2 := fd.LineContext(3)
	ctx2 := line2.Mark(15, 21)

	line3 := fd.LineContext(4)
	// 0         1         2         3         4
	// 0    5    0    5    0    5    0    5    0
	// ut labore et dolore magna aliqua
	ctx3 := line3.Mark(13, 19)

	ctx := Join(ctx1, ctx2, ctx3)
	ctx.Load(2, 2)

	got := ctx.HighlightText("the quick brown fox")
	expected := strings.Join([]string{
		"   2:   consectetur adipiscing elit",
		"   3:   ",
		"   4:   sed do eiusmod tempor incididunt",
		"               ^^^^^^^ ^^^^^^",
		"   5:   ut labore et dolore magna aliqua",
		"                     ^^^^^^",
		"                     the quick brown fox",
		"   6:   ut enim ad minim veniam",
		"   7:   ",
	}, "\n")

	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestContextHighlightTextMultipleParts3(t *testing.T) {
	fd := createTestFile1()

	line1 := fd.LineContext(3)
	// 0         1         2         3         4
	// 0    5    0    5    0    5    0    5    0
	// sed do eiusmod tempor incididunt
	ctx1 := line1.Mark(7, 14)

	line2 := fd.LineContext(3)
	ctx2 := line2.Mark(15, 21)

	line3 := fd.LineContext(4)
	// 0         1         2         3         4
	// 0    5    0    5    0    5    0    5    0
	// ut labore et dolore magna aliqua
	ctx3 := line3.Mark(13, 19)

	ctx := Join(ctx3, ctx2, ctx1)
	ctx.Load(2, 2)

	got := ctx.HighlightText("the quick brown fox")
	expected := strings.Join([]string{
		"   2:   consectetur adipiscing elit",
		"   3:   ",
		"   4:   sed do eiusmod tempor incididunt",
		"               ^^^^^^^ ^^^^^^",
		"   5:   ut labore et dolore magna aliqua",
		"                     ^^^^^^",
		"                     the quick brown fox",
		"   6:   ut enim ad minim veniam",
		"   7:   ",
	}, "\n")

	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestContextHighlightTextMultipleLines1(t *testing.T) {
	fd := createTestFile1()

	line1 := fd.LineContext(3)
	// 0         1         2         3         4
	// 0    5    0    5    0    5    0    5    0
	// sed do eiusmod tempor incididunt
	ctx1 := line1.Mark(7, 14)

	line2 := fd.LineContext(4)
	// 0         1         2         3         4
	// 0    5    0    5    0    5    0    5    0
	// ut labore et dolore magna aliqua
	ctx2 := line2.Mark(13, 19)

	ctx := Join(ctx1, ctx2)
	ctx.Load(2, 2)

	got := ctx.HighlightText("the quick brown fox")
	expected := strings.Join([]string{
		"   2:   consectetur adipiscing elit",
		"   3:   ",
		"   4:   sed do eiusmod tempor incididunt",
		"               ^^^^^^^",
		"   5:   ut labore et dolore magna aliqua",
		"                     ^^^^^^",
		"                     the quick brown fox",
		"   6:   ut enim ad minim veniam",
		"   7:   ",
	}, "\n")

	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestJoinContext(t *testing.T) {
	if c := Join(); c != nil {
		t.Fatalf("expected nil context, got non-nil")
	}

	fd := createTestFile1()

	line := fd.LineContext(3)
	// 0         1         2         3         4
	// 0    5    0    5    0    5    0    5    0
	// sed do eiusmod tempor incididunt
	ctx1 := line.Mark(7, 14)
	ctx1.Load(2, 2)

	ctx := Join(ctx1)

	got := ctx.HighlightText("the quick brown fox")
	expected := strings.Join([]string{
		"   2:   consectetur adipiscing elit",
		"   3:   ",
		"   4:   sed do eiusmod tempor incididunt",
		"               ^^^^^^^",
		"               the quick brown fox",
		"   5:   ut labore et dolore magna aliqua",
		"   6:   ut enim ad minim veniam",
	}, "\n")

	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}
