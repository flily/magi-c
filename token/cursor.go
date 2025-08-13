package token

import (
	"fmt"

	"github.com/flily/magi-c/context"
)

type Cursor struct {
	File   *context.FileContext
	Line   int
	Column int
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
	return c.Peek(0)
}

func (c *Cursor) Peek(n int) (rune, bool) {
	if c.Line >= len(c.File.Contents) {
		return 0, true
	}

	l := c.File.Line(c.Line)
	if c.Column+n >= len(l.Content) {
		return 0, true
	}

	return l.Content[c.Column+n], false
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

func (c *Cursor) Duplicate() *Cursor {
	cursor := &Cursor{
		File:   c.File,
		Line:   c.Line,
		Column: c.Column,
	}

	return cursor
}

func (c *Cursor) Start() *Cursor {
	return c.Duplicate()
}

func (c *Cursor) Finish(begin *Cursor) *context.Context {
	if begin.File != c.File {
		panic(fmt.Sprintf("cursor context file %s does not match cursor file %s", begin.File.Filename, c.File.Filename))
	}

	if begin.Line != c.Line {
		panic(fmt.Sprintf("cursor context line %d does not match cursor line %d", begin.Line, c.Line))
	}

	line := c.File.Line(begin.Line)
	return line.Mark(begin.Column, c.Column)
}
