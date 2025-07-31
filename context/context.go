package context

import (
	"fmt"
	"sort"
	"strings"
)

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
	PrevLines []*LineContext
	NextLines []*LineContext
	Lines     []*LineContext
}

func (c *Context) Join(ctxs ...*Context) *Context {
	lineMap := make(map[int]*LineContext)
	for _, line := range c.Lines {
		lineMap[line.Content.Line] = line.Duplicate()
	}

	for i, ctx := range ctxs {
		if ctx.File != c.File {
			err := fmt.Errorf("context %d does not match file", i)
			panic(err)
		}

		for _, l := range ctx.Lines {
			line, found := lineMap[l.Content.Line]
			if !found {
				lineMap[l.Content.Line] = l.Duplicate()
			} else {
				line.Highlights = append(line.Highlights, l.Highlights...)
			}
		}
	}

	result := &Context{
		File:  c.File,
		Lines: make([]*LineContext, 0, len(lineMap)),
	}

	for _, line := range lineMap {
		sort.Sort(ByHighlight(line.Highlights))
		result.Lines = append(result.Lines, line)
	}

	sort.Sort(ByLineContextLine(result.Lines))
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

	c.PrevLines = make([]*LineContext, 0, prev)
	c.NextLines = make([]*LineContext, 0, next)
	for i := 0; i < len(c.File.Contents); i++ {
		np, nl := lineFirst-i, i-lineLast
		if np > 0 && np <= prev {
			c.PrevLines = append(c.PrevLines, c.File.LineContext(i))
		}

		if nl > 0 && nl <= next {
			c.NextLines = append(c.NextLines, c.File.LineContext(i))
		}
	}
}

func (c *Context) HighlightTextWith(indicator string, format string, args ...any) string {
	parts := make([]string, 0, len(c.Lines)+len(c.PrevLines)+len(c.NextLines)+3)

	for _, line := range c.PrevLines {
		parts = append(parts, line.String())
	}

	for i, line := range c.Lines {
		if i == len(c.Lines)-1 {
			// show message in last line
			parts = append(parts, line.HighlighText(format, args...))
		} else {
			parts = append(parts, line.HighlighText(NoHighlightMessage))
		}
	}

	for _, line := range c.NextLines {
		parts = append(parts, line.String())
	}

	return strings.Join(parts, "\n")
}

func (c *Context) HighlightText(format string, args ...any) string {
	return c.HighlightTextWith(DefaultIndicator, format, args...)
}

func Join(ctxs ...*Context) *Context {
	if len(ctxs) == 0 {
		return nil
	}

	first := ctxs[0]
	if len(ctxs) == 1 {
		return first
	}

	return first.Join(ctxs[1:]...)
}
