package context

import (
	"fmt"
	"sort"
	"strings"
)

const (
	FixedLeadingSpace = "        "
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

func (l *LineContext) Mark(start int, end int) *LineContext {
	if start < 0 || end < 0 || start > end || start > l.Length() || end > l.Length() {
		err := fmt.Errorf("invalid context argument start=%d end=%d length=%d",
			start, end, l.Length())
		panic(err)
	}

	h := NewHighlight(start, end)
	l.Highlights = append(l.Highlights, h)
	sort.Sort(ByHighlight(l.Highlights))

	return l
}

func (l *LineContext) String() string {
	return fmt.Sprintf("%4d:   %s", l.Content.Line, l.Content.String())
}

func (l *LineContext) Highlight(format string, args ...any) string {
	message := fmt.Sprintf(format, args...)

	parts := make([]string, 0, 2*len(l.Highlights))
	last, lead := 0, ""
	for i, highlight := range l.Highlights {
		// highlight will store in order
		if highlight.Start < 0 || highlight.End > l.Length() || highlight.Start > highlight.End {
			err := fmt.Errorf("invalid highlight range: start=%d, end=%d, length=%d", highlight.Start, highlight.End, l.Length())
			panic(err)
		}

		widthSpace, widthHighligh := 0, 0
		for j := last; j < highlight.Start; j++ {
			widthSpace += CharWidthIn(l.Content.Content[j], j)
		}

		for j := highlight.Start; j < highlight.End; j++ {
			widthHighligh += CharWidthIn(l.Content.Content[j], j)
		}

		if i == 0 {
			// the first highlight
			lead = strings.Repeat(" ", widthSpace)
		}

		last = highlight.End
		parts = append(parts,
			strings.Repeat(" ", widthSpace),
			strings.Repeat("^", widthHighligh),
		)
	}

	content := l.String()
	return fmt.Sprintf("%s\n%s%s\n%s%s%s",
		content,
		FixedLeadingSpace, strings.Join(parts, ""),
		FixedLeadingSpace, strings.Repeat(" ", len(lead)), message,
	)
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

func (c *Context) Load(prev int, next int) {
	first := true
	lineFirst, lineLast := 100000000000, 0
	for _, line := range c.Lines {
		if first {
			lineFirst, lineLast = line.Content.Line, line.Content.Line

		} else {
			if line.Content.Line < lineFirst {
				lineFirst = line.Content.Line
			}

			if line.Content.Line > lineLast {
				lineLast = line.Content.Line
			}
		}

		first = false
	}

	c.PrevLines = make([]*LineContent, 0, prev)
	c.NextLines = make([]*LineContent, 0, next)
	for i := 0; i < len(c.File.Contents); i++ {
		np, nl := lineFirst-i, i-lineLast
		if np > 0 && np <= prev {
			c.PrevLines = append(c.PrevLines, c.File.Line(i))
		}

		if nl > 0 && nl <= next {
			c.NextLines = append(c.NextLines, c.File.Line(i))
		}
	}
}
