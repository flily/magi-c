package token

import (
	"testing"

	"bytes"
	"fmt"
	"strings"

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

	r1, eol := cursor.Peek(1)
	if eol {
		t.Fatalf("unexpected EOF at %s", cursor.Position())
	}

	if r1 != firstLine[1] {
		t.Errorf("expected rune '%c' at %s, got '%c'", firstLine[1], cursor.Position(), r1)
	}

	r2, eol := cursor.Peek(2)
	if eol {
		t.Fatalf("unexpected EOF at %s", cursor.Position())
	}

	if r2 != firstLine[2] {
		t.Errorf("expected rune '%c' at %s, got '%c'", firstLine[2], cursor.Position(), r2)
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

	rp, eol := cursor.Peek(1)
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

	eof := cursor.NextNonEmptyLine()
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
		eof := cursor.NextNonEmptyLine()
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

func TestCursorNextLineInThieLastLine(t *testing.T) {
	cursor := createTestCursor1()

	cursor.NextNonEmptyLine() // move to second line
	for range 5 {
		_, eof := cursor.NextInLine()
		if eof {
			t.Fatalf("unexpected EOF at %s", cursor.Position())
		}
	}

	r, eof := cursor.NextInLine()
	if !eof {
		t.Errorf("expected to be at end of line at %s, got rune '%c'", cursor.Position(), r)
	}

	eof = cursor.NextNonEmptyLine()
	if !eof {
		t.Errorf("expected to be at end of file at %s, got rune '%c'", cursor.Position(), r)
	}

	r, eof = cursor.NextInLine()
	if !eof {
		t.Errorf("expected to be at end of line at %s, got rune '%c'", cursor.Position(), r)
	}
}

func TestCursorMark(t *testing.T) {
	cursor := createTestCursor2()

	s := cursor.State()
	for {
		r, _ := cursor.NextInLine()
		if r == ' ' {
			break
		}
	}

	ctx := cursor.Finish(s)
	got := ctx.HighlightText("here")
	expected := strings.Join([]string{
		"   1:   lorem ipsum dolor sit amet",
		"        ^^^^^",
		"        here",
	}, "\n")

	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestCursorPeekString(t *testing.T) {
	cursor := createTestCursor2()

	expPosition := "example.txt:1:1"
	if cursor.Position() != expPosition {
		t.Errorf("expected position '%s', got '%s'", expPosition, cursor.Position())
	}

	ctx1 := cursor.PeekString("lorem")
	if ctx1 == nil {
		t.Fatalf("expected to find 'lorem' at %s", cursor.Position())
	}

	expected1 := strings.Join([]string{
		"   1:   lorem ipsum dolor sit amet",
		"        ^^^^^",
		"        here",
	}, "\n")
	if expected1 != ctx1.HighlightText("here") {
		t.Errorf("expected:\n%s\ngot:\n%s", expected1, ctx1.HighlightText("here"))
	}

	if ctx := cursor.PeekString("ipsum"); ctx != nil {
		t.Errorf("expected to not find 'ipsum' at %s, got %s", cursor.Position(), ctx.HighlightText("here"))
	}
}

func TestCursorNextLiteral(t *testing.T) {
	cursor := createTestCursor1()

	expPosition := "example.txt:1:1"
	if cursor.Position() != expPosition {
		t.Errorf("expected position '%s', got '%s'", expPosition, cursor.Position())
	}

	ctx1 := cursor.PeekString("lorem")
	if ctx1 == nil {
		t.Fatalf("expected to find 'lorem' at %s", cursor.Position())
	}

	expected1 := strings.Join([]string{
		"   1:   lorem",
		"        ^^^^^",
		"        here",
	}, "\n")
	if expected1 != ctx1.HighlightText("here") {
		t.Errorf("expected:\n%s\ngot:\n%s", expected1, ctx1.HighlightText("here"))
	}

	expPosition = "example.txt:1:1"
	if cursor.Position() != expPosition {
		t.Errorf("expected position '%s', got '%s'", expPosition, cursor.Position())
	}

	ctx2 := cursor.NextString("ipsum")
	if ctx2 != nil {
		t.Fatalf("expected to not find 'ipsum' at %s, got %s", cursor.Position(), ctx2.HighlightText("here"))
	}

	ctx3 := cursor.NextString("lorem")
	if expected1 != ctx3.HighlightText("here") {
		t.Errorf("expected:\n%s\ngot:\n%s", expected1, ctx3.HighlightText("here"))
	}

	expPosition = "example.txt:1:6"
	if cursor.Position() != expPosition {
		t.Errorf("expected position '%s', got '%s'", expPosition, cursor.Position())
	}
}

func TestCursorNextLiteralInEndOfLine(t *testing.T) {
	cursor := createTestCursor1()

	ctx1 := cursor.NextString("lorem")
	if ctx1 == nil {
		t.Fatalf("expected to find 'lorem' at %s", cursor.Position())
	}

	expected1 := strings.Join([]string{
		"   1:   lorem",
		"        ^^^^^",
		"        here",
	}, "\n")
	if expected1 != ctx1.HighlightText("here") {
		t.Errorf("expected:\n%s\ngot:\n%s", expected1, ctx1.HighlightText("here"))
	}

	cursor.NextNonEmptyLine()

	ctx2 := cursor.NextString("ipsum")
	if ctx2 == nil {
		t.Fatalf("expected to find 'ipsum' at %s", cursor.Position())
	}
	expected2 := strings.Join([]string{
		"   2:   ipsum",
		"        ^^^^^",
		"        here",
	}, "\n")
	if expected2 != ctx2.HighlightText("here") {
		t.Errorf("expected:\n%s\ngot:\n%s", expected2, ctx2.HighlightText("here"))
	}
}

func TestCursorNext(t *testing.T) {
	cursor := createTestCursor1()

	first, _ := cursor.Rune()
	result := make([]rune, 0, 100)
	result = append(result, first)
	for {
		r, eof := cursor.Next()
		if eof {
			break
		}

		result = append(result, r)
	}

	if len(result) != 10 {
		t.Errorf("expected 10 runes, got %d", len(result))
		t.Errorf("result: [%s]", string(result))
	}
}
