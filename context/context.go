package context

import (
	"fmt"
	"sort"
)

type Highlight struct {
	Start int
	End   int
}

func NewHighlight(start, end int) Highlight {
	if start < 0 || end < 0 || start > end {
		err := fmt.Errorf("invalid highlight range: start=%d, end=%d", start, end)
		panic(err)
	}

	h := Highlight{
		Start: start,
		End:   end,
	}

	return h
}

type ByHighlight []Highlight

func (a ByHighlight) Len() int {
	return len(a)
}

func (a ByHighlight) Less(i, j int) bool {
	return a[i].Start < a[j].Start
}

func (a ByHighlight) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type Line struct {
	Line    int
	Content []rune
}

type LineContext struct {
	Content    *LineContent
	File       *FileContext
	Highlights []Highlight
}

func (l *LineContext) StringContent() string {
	return l.Content.String()
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

func (l *LineContext) Join(lctxs ...*LineContext) *LineContext {
	hc := len(l.Highlights)
	for _, lctx := range lctxs {
		hc += len(lctx.Highlights)
	}

	result := &LineContext{
		Content:    l.Content,
		File:       l.File,
		Highlights: make([]Highlight, 0, hc),
	}

	result.Highlights = append(result.Highlights, l.Highlights...)

	for _, lctx := range lctxs {
		if lctx.Content != l.Content {
			err := fmt.Errorf("cannot join different line contexts: %s != %s", l.Content.String(), lctx.Content.String())
			panic(err)
		}

		result.Highlights = append(result.Highlights, lctx.Highlights...)
	}

	sort.Sort(ByHighlight(l.Highlights))
	return result
}

type ByLineContextLine []*LineContext

func (a ByLineContextLine) Len() int {
	return len(a)
}

func (a ByLineContextLine) Less(i, j int) bool {
	return a[i].Content.Line < a[j].Content.Line
}

func (a ByLineContextLine) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type Context struct {
	File      *FileContext
	PrevLines []*LineContent
	NextLines []*LineContent
	Lines     []*LineContext
}

func (c *Context) Join(ctx *Context) *Context {
	if c.File != ctx.File {
		return nil
	}

	result := &Context{
		File: c.File,
	}

	for _, line := range c.Lines {
		for _, l := range ctx.Lines {
			if line.Content == l.Content {
				result.Lines = append(result.Lines, line.Join(l))
			}
		}
	}

	return result
}
