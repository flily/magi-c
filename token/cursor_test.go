package token

import (
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

func TestCursorBasic1(t *testing.T) {
	cursor := createTestCursor1()

	firstLine := []rune{
		'l', 'o', 'r', 'e', 'm',
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

	for range 5 {
		_, eol := cursor.NextInLine()
		if !eol {
			t.Errorf("expected to be at end of line at %s", cursor.Position())
		}
	}
}
