package context

import (
	"testing"

	"bytes"
)

func TestLineContentBasicInfo(t *testing.T) {
	s := "lorem ipsum dolor sit amet"
	line := NewLineFromBytes(42, []byte(s), nil)
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
	line := NewLineFromBytes(42, []byte(s), nil)
	ctx := line.MarkLine(6, 11)

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

	ctx := ReadFileData(filename, bytes.Join(content, nil))
	if ctx.Filename != filename {
		t.Errorf("expected filename %s, got %s", filename, ctx.Filename)
	}

	if ctx.Lines() != 7 {
		t.Errorf("expected 7 lines, got %d", ctx.Lines())
	}

	if ctx.Line(0).String() != "lorem ipsum" {
		t.Errorf("expected first line to be 'lorem ipsum', got '%s'", ctx.Line(0).String())
	}

	if !bytes.Equal(ctx.Line(0).EOL, EolLF) {
		t.Errorf("expected first line EOL to be LF, got %s", ctx.Line(0).EOL)
	}

	if ctx.Line(1).String() != "dolor sit amet" {
		t.Errorf("expected second line to be 'dolor sit amet', got '%s'", ctx.Line(1).String())
	}

	if !bytes.Equal(ctx.Line(1).EOL, EolCRLF) {
		t.Errorf("expected second line EOL to be CRLF, got %s", ctx.Line(1).EOL)
	}

	if ctx.Line(6).String() != "ut labore et dolore magna aliqua" {
		t.Errorf("expected last line to be 'ut labore et dolore magna aliqua', got '%s'", ctx.Line(6).String())
	}

	if c, eol := ctx.Rune(0, 0); c != 'l' || eol {
		t.Errorf("expected first rune to be 'l', got '%c', eof=%v", c, eol)
	}

	if c, eol := ctx.Rune(4, 5); c != 'o' || eol {
		t.Errorf("expected rune at (4, 5) to be 'o', got '%c', eof=%v", c, eol)
	}

	if c, eol := ctx.Rune(0, 100); c != 0 || !eol {
		t.Errorf("expected rune at (0, 100) to be EOF, got '%c', eof=%v", c, eol)
	}

	if c, eol := ctx.Rune(100, 0); c != 0 || !eol {
		t.Errorf("expected rune at (100, 0) to be EOF, got '%c', eof=%v", c, eol)
	}
}
