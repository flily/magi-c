package context

import (
	"testing"
)

func TestLineContentBasicInfo(t *testing.T) {
	s := "lorem ipsum dolor sit amet"
	line := NewLineFromBytes(42, []byte(s))
	if line.Line != 42 {
		t.Errorf("expected line number 42, got %d", line.Line)
	}

	if line.String() != s {
		t.Errorf("expected line content '%s', got '%s'", s, line.String())
	}

	if line.Length() != len(s) {
		t.Errorf("expected line length %d, got %d", len(s), line.Length())
	}
}

func TestLineContextMark(t *testing.T) {
	s := "lorem ipsum dolor sit amet"
	line := NewLineFromBytes(42, []byte(s))
	ctx := line.Mark(6, 11)

	if ctx.StringContent() != s {
		t.Errorf("expected context content '%s', got '%s'", s, ctx.StringContent())
	}

	if ctx.Length() != len(s) {
		t.Errorf("expected context length %d, got %d", len(s), ctx.Length())
	}

	if len(ctx.Highlights) != 1 {
		t.Errorf("expected 1 highlight, got %d", len(ctx.Highlights))
	}
}

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

	if ctx.LineContent(0).String() != "lorem ipsum" {
		t.Errorf("expected first line to be 'lorem ipsum', got '%s'", ctx.LineContent(0).String())
	}

	if ctx.LineContent(6).String() != "ut labore et dolore magna aliqua" {
		t.Errorf("expected last line to be 'ut labore et dolore magna aliqua', got '%s'", ctx.LineContent(6).String())
	}
}
