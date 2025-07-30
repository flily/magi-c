package context

import (
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

func (c *Context) Join(ctx *Context) *Context {
	if c.File != ctx.File {
		return nil
	}

	result := &Context{
		File: c.File,
	}

	for _, line := range c.Lines {
		lctx := FindLineContextSameLine(ctx.Lines, line)
		if lctx == nil {
			result.Lines = append(result.Lines, line)
		} else {
			result.Lines = append(result.Lines, line.Join(lctx))
		}
	}

	for _, line := range ctx.Lines {
		if lctx := FindLineContextSameLine(c.Lines, line); lctx == nil {
			result.Lines = append(result.Lines, line)
		}
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
