package context

import (
	"testing"

	"bytes"
	"fmt"
	"strings"
)

func makeTestCursor(filename string, data []byte) *Cursor {
	file := ReadFileData(filename, data)
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
		r, eol, eof := cursor.Rune()
		if eol || eof {
			t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.Position(), eol, eof)
		}

		if r != expChar {
			t.Errorf("expected rune '%c' at %s, got '%c'", expChar, cursor.Position(), r)
		}

		expPosition := fmt.Sprintf("example.txt:1:%d", i+1)
		if cursor.Position() != expPosition {
			t.Errorf("expected position '%s', got '%s'", expPosition, cursor.Position())
		}

		_, eol = cursor.NextInLine()
		if i != len(firstLine)-1 && eol {
			t.Errorf("unexpected EOL at %s", cursor.Position())
		}
	}

	for range 5 {
		_, eol := cursor.NextInLine()
		if !eol {
			t.Errorf("expected EOL at %s", cursor.Position())
		}
	}
}

func TestCursorCurrentChar(t *testing.T) {
	cursor := createTestCursor1()

	c1 := cursor.CurrentChar()
	exp1 := strings.Join([]string{
		"   1:   lorem",
		"        ^",
		"        here",
	}, "\n")
	got1 := c1.HighlightText("here")
	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	_, _, _ = cursor.Next()
	c2 := cursor.CurrentChar()
	exp2 := strings.Join([]string{
		"   1:   lorem",
		"         ^",
		"         here",
	}, "\n")
	got2 := c2.HighlightText("here")
	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	_ = cursor.NextNonEmptyLine()
	c3 := cursor.CurrentChar()
	exp3 := strings.Join([]string{
		"   2:   ipsum",
		"        ^",
		"        here",
	}, "\n")
	got3 := c3.HighlightText("here")
	if got3 != exp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp3, got3)
	}
}

func TestCursorPeek(t *testing.T) {
	cursor := createTestCursor1()

	firstLine := []rune{
		'l', 'o', 'r', 'e', 'm',
	}

	r0, eol, eof := cursor.Rune()
	if eol || eof {
		t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.Position(), eol, eof)
	}

	if r0 != firstLine[0] {
		t.Errorf("expected rune '%c' at %s, got '%c'", firstLine[0], cursor.Position(), r0)
	}

	r1, eol, eof := cursor.Peek(1)
	if eol || eof {
		t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.Position(), eol, eof)
	}

	if r1 != firstLine[1] {
		t.Errorf("expected rune '%c' at %s, got '%c'", firstLine[1], cursor.Position(), r1)
	}

	r2, eol, eof := cursor.Peek(2)
	if eol || eof {
		t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.Position(), eol, eof)
	}

	if r2 != firstLine[2] {
		t.Errorf("expected rune '%c' at %s, got '%c'", firstLine[2], cursor.Position(), r2)
	}

	// move to end of line
	for range 4 {
		cursor.NextInLine()
	}

	re, eol, eof := cursor.Rune()
	if eol || eof {
		t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.Position(), eol, eof)
	}

	if re != firstLine[4] {
		t.Errorf("expected rune '%c' at %s, got '%c'", firstLine[4], cursor.Position(), re)
	}

	rp, eol, eof := cursor.Peek(1)
	if !eol || eof {
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

	for i, expChar := range firstLine {
		r, eol, eof := cursor.Rune()
		if eol || eof {
			t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.Position(), eol, eof)
		}

		if r != expChar {
			t.Errorf("expected rune '%c' at %s, got '%c'", expChar, cursor.Position(), r)
		}

		_, eol = cursor.NextInLine()
		if i != len(firstLine)-1 && eol {
			t.Errorf("unexpected EOL at %s", cursor.Position())
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
		r, eol, eof := cursor.Rune()
		if eol || eof {
			t.Fatalf("unexpected EOL/EOF at %s  (EOL=%v, EOF=%v)", cursor.Position(), eol, eof)
		}

		if r != expChar {
			t.Errorf("expect rune '%c' at %s, got '%c'", expChar, cursor.Position(), r)
		}

		_, eol = cursor.NextInLine()
		if eol {
			t.Errorf("unexpected EOL at %s", cursor.Position())
		}
	}
}

func TestCursorNextLineInTheLastLine(t *testing.T) {
	cursor := createTestCursor1()

	cursor.NextNonEmptyLine() // move to second line
	for range 4 {
		_, eof := cursor.NextInLine()
		if eof {
			t.Fatalf("unexpected EOF at %s", cursor.Position())
		}
	}

	// 5th char
	_, eof := cursor.NextInLine()
	if !eof {
		t.Errorf("expect EOL at %s", cursor.Position())
	}

	r, eof := cursor.NextInLine()
	if !eof {
		t.Errorf("expect EOL at %s, got rune '%c'", cursor.Position(), r)
	}

	eof = cursor.NextNonEmptyLine()
	if !eof {
		t.Errorf("expect EOL at %s, got rune '%c'", cursor.Position(), r)
	}

	r, eof = cursor.NextInLine()
	if !eof {
		t.Errorf("expect EOL at %s, got rune '%c'", cursor.Position(), r)
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

	content, ctx := cursor.Finish(s)
	got := ctx.HighlightText("here")
	expected := strings.Join([]string{
		"   1:   lorem ipsum dolor sit amet",
		"        ^^^^^",
		"        here",
	}, "\n")

	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}

	if content != "lorem" {
		t.Errorf("expected content 'lorem', got '%s'", content)
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

	first, _, _ := cursor.Rune()
	result := make([]rune, 0, 100)
	result = append(result, first)
	for {
		r, _, eof := cursor.Next()
		if eof {
			break
		}

		result = append(result, r)
	}

	if len(result) != 12 {
		t.Errorf("expected 12 runes, got %d", len(result))
		t.Errorf("result: [%s]", string(result))
	}
}

func TestCursorIsFirstNonWhiteChar2(t *testing.T) {
	text := [][]byte{
		[]byte("the quick brown fox\n"),
		[]byte("   jumps over the lazy dog\n"),
	}

	cursor := makeTestCursor("example.txt", bytes.Join(text, []byte{}))

	if !cursor.IsFirstNonWhiteChar() {
		t.Errorf("expected to be first non-white char at %s", cursor.Position())
	}

	for {
		_, eol := cursor.NextInLine()
		if eol {
			break
		}

		if cursor.IsFirstNonWhiteChar() {
			t.Errorf("expected to not be first non-white char at %s", cursor.Position())
		}
	}

	cursor.NextNonEmptyLine()

	if !cursor.IsFirstNonWhiteChar() {
		t.Errorf("expected to be first non-white char at %s", cursor.Position())
	}

	for {
		r, eol := cursor.NextInLine()
		if eol {
			break
		}

		if !cursor.IsFirstNonWhiteChar() {
			t.Errorf("expected to be first non-white char at %s", cursor.Position())
		}

		if r != ' ' {
			break
		}
	}

	for {
		_, eol := cursor.NextInLine()
		if eol {
			break
		}

		if cursor.IsFirstNonWhiteChar() {
			t.Errorf("expected to not be first non-white char at %s", cursor.Position())
		}
	}
}

func TestCursorSkipWhitespace(t *testing.T) {
	text := []byte(strings.Join([]string{
		"                lorem ipsum",
	}, "\n"))

	cursor := makeTestCursor("example.txt", text)

	got1 := cursor.CurrentChar().HighlightText("here")
	exp1 := strings.Join([]string{
		"   1:                   lorem ipsum",
		"        ^",
		"        here",
	}, "\n")

	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	cursor.SkipWhitespace()

	got2 := cursor.CurrentChar().HighlightText("here")
	exp2 := strings.Join([]string{
		"   1:                   lorem ipsum",
		"                        ^",
		"                        here",
	}, "\n")

	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}
}

func TestCursorSkipWhitespaceToNextLine(t *testing.T) {
	text := []byte(strings.Join([]string{
		"        lorem        ",
		"    ",
		"        \t\t\t        ",
		"      ipsum dolor sit amet",
	}, "\n"))

	cursor := makeTestCursor("example.txt", text)

	got1 := cursor.CurrentChar().HighlightText("here")
	exp1 := strings.Join([]string{
		"   1:           lorem        ",
		"        ^",
		"        here",
	}, "\n")

	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	cursor.SkipWhitespace()

	got2 := cursor.CurrentChar().HighlightText("here")
	exp2 := strings.Join([]string{
		"   1:           lorem        ",
		"                ^",
		"                here",
	}, "\n")

	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	cursor.Skip(5)
	cursor.SkipWhitespace()

	got3 := cursor.CurrentChar().HighlightText("here")
	exp3 := strings.Join([]string{
		"   4:         ipsum dolor sit amet",
		"              ^",
		"              here",
	}, "\n")

	if got3 != exp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp3, got3)
	}
}

func TestCursorSkipWhitespaceInLine(t *testing.T) {
	text := []byte(strings.Join([]string{
		"        lorem        ",
		"    ",
		"        \t\t\t        ",
		"      ipsum dolor sit amet",
	}, "\n"))

	cursor := makeTestCursor("example.txt", text)

	got1 := cursor.CurrentChar().HighlightText("here")
	exp1 := strings.Join([]string{
		"   1:           lorem        ",
		"        ^",
		"        here",
	}, "\n")

	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	cursor.SkipWhitespaceInLine()

	got2 := cursor.CurrentChar().HighlightText("here")
	exp2 := strings.Join([]string{
		"   1:           lorem        ",
		"                ^",
		"                here",
	}, "\n")

	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	cursor.Skip(5)
	cursor.SkipWhitespaceInLine()

	got3 := cursor.CurrentChar().HighlightText("here")
	exp3 := strings.Join([]string{
		"   1:           lorem        <EOL LF>",
		"                             ^^^^^^^^",
		"                             here",
	}, "\n")

	if got3 != exp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp3, got3)
	}
}
