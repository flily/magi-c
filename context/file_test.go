package context

import (
	"testing"
)

func TestReadFileData(t *testing.T) {
	filename := "example.txt"
	content := [][]byte{
		[]byte("lorem ipsum\n"),
		[]byte("dolor sit amet\r\n"),
		[]byte("consectetur adipiscing elit\n"),
		[]byte("\n"),
		[]byte("sed do eiusmod tempor incididunt\r\n"),
		[]byte("\r\n"),
		[]byte("ut labore et dolore magna aliqua\n"),
	}

	data := make([]byte, 0, 1024)
	for _, line := range content {
		data = append(data, line...)
	}

	ctx := ReadFileData(filename, data)
	if ctx.Filename != filename {
		t.Errorf("expected filename %s, got %s", filename, ctx.Filename)
	}

	if ctx.Lines() != 7 {
		t.Errorf("expected 7 lines, got %d", ctx.Lines())
	}

	if ctx.Line(0).String() != "lorem ipsum" {
		t.Errorf("expected first line to be 'lorem ipsum', got '%s'", ctx.Line(0).String())
	}

	if ctx.Line(6).String() != "ut labore et dolore magna aliqua" {
		t.Errorf("expected last line to be 'ut labore et dolore magna aliqua', got '%s'", ctx.Line(6).String())
	}
}
