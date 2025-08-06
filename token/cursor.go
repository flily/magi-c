package token

import (
	"fmt"

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

func (c *Cursor) Position() string {
	return fmt.Sprintf("%s:%d:%d", c.File.Filename, c.Line+1, c.Column+1)
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
	if c.Column+1 >= len(l.Content) {
		return 0, true
	}

	return l.Content[c.Column+1], false
}

func (c *Cursor) next(line int, column int) (int, int, *context.LineContent) {
	content := c.File.Line(line)
	if content == nil {
		return line, column, nil
	}

	column += 1
	if column >= content.Length() {
		c.NextLine()
		column = 0
	}

	return line, column, content
}

func (c *Cursor) nextInLine(line int, column int) (int, int, *context.LineContent) {
	content := c.File.Line(line)
	if content == nil {
		return line, column, nil
	}

	if c.Column >= content.Length() {
		return line, column, nil
	}

	return line, column + 1, content
}

func (c *Cursor) NextInLine() (rune, bool) {
	line, column, content := c.nextInLine(c.Line, c.Column)
	if content == nil {
		return 0, true
	}

	r := content.Rune(column) // checked in nextInLine()
	c.Line = line
	c.Column = column
	return r, false
}

func (c *Cursor) NextLine() bool {
	line := c.Line + 1
	content := c.File.Line(line)
	for content != nil && content.Length() <= 0 {
		line += 1
		content = c.File.Line(line)
	}

	c.Line = line
	c.Column = 0
	return content == nil
}

func (c *Cursor) Next() (rune, bool) {
	line, column, content := c.next(c.Line, c.Column)
	if content == nil {
		return 0, true
	}

	c.Line = line
	c.Column = column
	return c.Rune()
}
