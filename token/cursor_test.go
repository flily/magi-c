package token

import (
	"fmt"
	"testing"

	"bytes"

	"github.com/flily/magi-c/context"
)

func makeTestCursor(filename string, data []byte) *Cursor {
	file := context.ReadFileData(filename, data)
	cursor := NewCursor(file)
	return cursor
}

func createTestCursor1() *Cursor {
	parts := [][]byte{
		[]byte("lorem\n"),
		[]byte("ipsum\n"),
	}

	return makeTestCursor("example.txt", bytes.Join(parts, []byte{}))
}

func createTestCursor2() *Cursor {
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

	return makeTestCursor("example.txt", bytes.Join(parts, []byte{}))
}

func TestCursorRuneBasic(t *testing.T) {
	cursor := createTestCursor1()

	firstLine := []rune{
		'l', 'o', 'r', 'e', 'm',
	}

	for i, expChar := range firstLine {
		r, eof := cursor.Rune()
		if eof {
			t.Fatalf("unexpected EOF at %s", cursor.Position())
		}

		if r != expChar {
			t.Errorf("expected rune '%c' at %s, got '%c'", expChar, cursor.Position(), r)
		}

		expPosition := fmt.Sprintf("example.txt:1:%d", i+1)
		if cursor.Position() != expPosition {
			t.Errorf("expected position '%s', got '%s'", expPosition, cursor.Position())
		}

		_, eol := cursor.NextInLine()
		if eol {
			t.Errorf("expected not to be at end of line at %s", cursor.Position())
		}
	}

	for range 5 {
		_, eol := cursor.NextInLine()
		if !eol {
			t.Errorf("expected to be at end of line at %s", cursor.Position())
		}
	}
}

func TestCursorPeek(t *testing.T) {
	cursor := createTestCursor1()

	firstLine := []rune{
		'l', 'o', 'r', 'e', 'm',
	}

	r0, eol := cursor.Rune()
	if eol {
		t.Fatalf("unexpected EOF at %s", cursor.Position())
	}

	if r0 != firstLine[0] {
		t.Errorf("expected rune '%c' at %s, got '%c'", firstLine[0], cursor.Position(), r0)
	}

	r1, eol := cursor.Peek()
	if eol {
		t.Fatalf("unexpected EOF at %s", cursor.Position())
	}

	if r1 != firstLine[1] {
		t.Errorf("expected rune '%c' at %s, got '%c'", firstLine[1], cursor.Position(), r1)
	}

	// move to end of line
	for range 4 {
		cursor.NextInLine()
	}

	re, eol := cursor.Rune()
	if eol {
		t.Fatalf("unexpected EOF at %s", cursor.Position())
	}

	if re != firstLine[4] {
		t.Errorf("expected rune '%c' at %s, got '%c'", firstLine[4], cursor.Position(), re)
	}

	rp, eol := cursor.Peek()
	if !eol {
		t.Errorf("expected to be at end of line at %s, got rune '%c'", cursor.Position(), rp)
	}

	if rp != 0 {
		t.Errorf("expected rune to be 0 at end of line, got '%c'", rp)
	}
}

func TestCursorNextLineBasic(t *testing.T) {
	cursor := createTestCursor1()

	firstLine := []rune{
		'i', 'p', 's', 'u', 'm',
	}

	eof := cursor.NextLine()
	if eof {
		t.Fatalf("unexpected EOF at %s", cursor.Position())
	}

	for _, expChar := range firstLine {
		r, eof := cursor.Rune()
		if eof {
			t.Fatalf("unexpected EOF at %s", cursor.Position())
		}

		if r != expChar {
			t.Errorf("expected rune '%c' at %s, got '%c'", expChar, cursor.Position(), r)
		}

		_, eol := cursor.NextInLine()
		if eol {
			t.Errorf("expected not to be at end of line at %s", cursor.Position())
		}
	}
}

func TestCursorNextLineWithEmptyLine(t *testing.T) {
	cursor := createTestCursor2()

	content := []rune{
		's', 'e', 'd',
	}

	for range 2 {
		eof := cursor.NextLine()
		if eof {
			t.Fatalf("unexpected EOF at %s", cursor.Position())
		}
	}

	for _, expChar := range content {
		r, eof := cursor.Rune()
		if eof {
			t.Fatalf("unexpected EOF at %s", cursor.Position())
		}

		if r != expChar {
			t.Errorf("expected rune '%c' at %s, got '%c'", expChar, cursor.Position(), r)
		}

		_, eol := cursor.NextInLine()
		if eol {
			t.Errorf("expected not to be at end of line at %s", cursor.Position())
		}
	}
}
