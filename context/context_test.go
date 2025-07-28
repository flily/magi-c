package context

import (
	"bytes"
	"strings"
	"testing"
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
