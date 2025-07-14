package context

type Cursor struct {
	Line        int
	Column      int
	File        *FileContext
	LineContext *LineContext
}

func (c *Cursor) Next() *Cursor {
	lctx := c.LineContext
	line := c.Line
	column := c.Column + 1
	for column >= c.LineContext.Length() {
		line += 1
		lctx = c.File.LineContext(line)
		if lctx == nil {
			return nil
		}
		column = 0
	}

	return &Cursor{
		Line:        line,
		Column:      column,
		File:        c.File,
		LineContext: lctx,
	}
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
