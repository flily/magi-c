package token

import (
	"github.com/flily/magi-c/context"
)

func IsValidIdentifierRune(r rune) bool {
	if 'a' <= r && r <= 'z' {
		return true
	}

	if 'A' <= r && r <= 'Z' {
		return true
	}

	if '0' <= r && r <= '9' {
		return true
	}

	if r == '_' {
		return true
	}

	return false
}

func IsValidIdentifierInitialRune(r rune) bool {
	if 'a' <= r && r <= 'z' {
		return true
	}

	if 'A' <= r && r <= 'Z' {
		return true
	}

	if r == '_' {
		return true
	}

	return false
}

type Cursor struct {
	Line   int
	Column int
	File   *context.FileContext
}

func NewCursor(file *context.FileContext) *Cursor {
	c := &Cursor{
		File: file,
	}

	return c
}

func (c *Cursor) Rune() (rune, bool) {
	if c.Line >= len(c.File.Contents) {
		return 0, true
	}

	// Cursor will never stop on a column that is out of bounds
	l := c.File.Line(c.Line)
	return l.Content[c.Column], false
}

func (c *Cursor) Peek() (rune, bool) {
	if c.Line >= len(c.File.Contents) {
		return 0, true
	}

	l := c.File.Line(c.Line)
	if c.Column >= len(l.Content) {
		return 0, true
	}

	return l.Content[c.Column], false
}

func (c *Cursor) moveNext(line int, column int) (int, int, *context.LineContent) {
	content := c.File.Line(line)
	if content == nil {
		return line, column, nil
	}

	column += 1
	for column >= content.Length() {
		line += 1
		content = c.File.Line(line)
		if content == nil {
			return line, 0, nil
		}
		column = 0
	}

	return line, column, content
}

func (c *Cursor) Next() (rune, bool) {
	line, column, content := c.moveNext(c.Line, c.Column)
	if content == nil {
		return 0, true
	}

	r := content.Rune(column) // checked in moveNext()
	c.Line = line
	c.Column = column
	return r, false
}
