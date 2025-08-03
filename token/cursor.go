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
	if c.Column >= len(l.Content) {
		return 0, true
	}

	return l.Content[c.Column], false
}

func (c *Cursor) nextInLine() (int, int, *context.LineContent) {
	if c.Line >= len(c.File.Contents) {
		return c.Line, c.Column, nil
	}

	content := c.File.Line(c.Line)
	if content == nil {
		return c.Line, c.Column, nil
	}

	if c.Column >= content.Length() {
		return c.Line, c.Column, nil
	}

	return c.Line, c.Column + 1, content
}

func (c *Cursor) NextInLine() (rune, bool) {
	line, column, content := c.nextInLine()
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
	}

	return content == nil
}

func (c *Cursor) next() *context.LineContent {
	content := c.File.Line(c.Line)
	if content == nil {
		return nil
	}

	c.Column += 1
	if c.Column >= content.Length() {
		c.NextLine()
	}

	return content
}

func (c *Cursor) Next() (rune, bool) {
	content := c.next()
	if content == nil {
		return 0, true
	}

	return c.Rune()
}
