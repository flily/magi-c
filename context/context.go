package context

import (
	"fmt"
)

type Highlight struct {
	Start int
	End   int
}

func NewHighlight(start, end int) *Highlight {
	if start < 0 || end < 0 || start > end {
		err := fmt.Errorf("invalid highlight range: start=%d, end=%d", start, end)
		panic(err)
	}

	h := &Highlight{
		Start: start,
		End:   end,
	}

	return h
}

type LineContext struct {
	Content    *LineContent
	File       *FileContext
	Highlights []*Highlight
}

func (l *LineContext) Length() int {
	return l.Content.Length()
}

func (l *LineContext) Rune(n int) (rune, bool) {
	if n < 0 || n >= l.Length() {
		return 0, true
	}

	r := l.Content.Content[n]
	return r, false
}

type Context struct {
	File      *FileContext
	PrevLines []*LineContent
	NextLines []*LineContent
	Line      *LineContext
}
