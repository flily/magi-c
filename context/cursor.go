package context

type Cursor struct {
	Line        int
	Column      int
	File        *FileContext
	LineContext *LineContext
}

func (c *Cursor) Rune() (rune, bool) {
	if c.Line >= len(c.File.Contents) {
		return 0, true
	}

	// Cursor will never stop on a column that is out of bounds
	l := c.File.LineContent(c.Line)
	return l.Content[c.Column], false
}

func (c *Cursor) moveNext(line int, column int) (int, int, *LineContext) {
	lctx := c.LineContext
	column += 1
	for column >= lctx.Length() {
		line += 1
		lctx = c.File.LineContext(line)
		if lctx == nil {
			return line, 0, nil
		}
		column = 0
	}

	return line, column, lctx
}

func (c *Cursor) Peek() (rune, bool) {
	if c.Line >= len(c.File.Contents) {
		return 0, true
	}

	l := c.File.LineContent(c.Line)
	if c.Column >= len(l.Content) {
		return 0, true
	}

	return l.Content[c.Column], false
}

func (c *Cursor) Next() (rune, bool) {
	line, column, lctx := c.moveNext(c.Line, c.Column)
	if lctx == nil {
		return 0, true
	}

	r := lctx.Content.Content[column]
	c.LineContext = lctx
	c.Line = line
	c.Column = column
	return r, false
}

func (c *Cursor) NextCursor() *Cursor {
	next := &Cursor{
		Line:        c.Line,
		Column:      c.Column,
		File:        c.File,
		LineContext: c.LineContext,
	}

	_, eof := next.Next()
	if eof {
		return nil
	}

	return next
}
