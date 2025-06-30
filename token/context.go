package token

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"
)

const (
	NoLineNumber = -1
)

// HighlightContext represents the start (included) and end (excluded) position of the highlight.
type HighlightContext struct {
	Start int
	End   int
}

type ByHighlight []HighlightContext

func (a ByHighlight) Len() int {
	return len(a)
}

func (a ByHighlight) Less(i, j int) bool {
	return a[i].Start < a[j].Start
}

func (a ByHighlight) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type Line struct {
	Line    int
	Content []rune
}

func NewLineFromBytes(line int, data []byte) Line {
	l := Line{
		Line:    line,
		Content: []rune(string(data)),
	}

	return l
}

func (l Line) Valid() bool {
	return l.Line >= 0
}

func (l Line) Length() int {
	return len(l.Content)
}

func (l Line) AsString() string {
	return string(l.Content)
}

func (l Line) AsRunes() []rune {
	return l.Content
}

// Char returns the character at the specified position, and end of line flag.
func (l Line) Char(n int) (rune, bool) {
	if n < 0 || n >= l.Length() {
		return 0, true
	}

	return l.Content[n], false

}

func NewInvalidLine() Line {
	return Line{Line: -1}
}

type LineContext struct {
	Filename   string
	Content    Line
	Highlights []HighlightContext
}

func NewLineContext(filename string, content Line) *LineContext {
	ctx := &LineContext{
		Filename:   filename,
		Content:    content,
		Highlights: make([]HighlightContext, 0, 4),
	}

	return ctx
}

func NewLineContextFromString(filename string, line int, content string) *LineContext {
	l := NewLineFromBytes(line, []byte(content))
	return NewLineContext(filename, l)
}

func NewEmptyLineContext(filename string, line int) *LineContext {
	return NewLineContext(filename, NewLineFromBytes(line, []byte{}))
}

func (c *LineContext) AddHighlight(start int, end int) {
	highlight := HighlightContext{
		Start: start,
		End:   end,
	}

	c.Highlights = append(c.Highlights, highlight)
}

func (c *LineContext) MakeLineWithNumber(format string) string {
	line := c.Content.AsString()

	if len(format) > 0 {
		// Line number in memory starts from 0, but we want to display it starting from 1.
		number := fmt.Sprintf(format, c.Content.Line+1)
		return number + line

	} else {
		return line
	}
}

func (c *LineContext) MakeHighlight(placeholder string, message string) (string, string) {
	lead := 0
	hasHighlighted := false
	lastCh := 0
	for _, h := range c.Highlights {
		if h.End > lastCh {
			lastCh = h.End
		}
	}

	if lastCh <= 0 || len(c.Highlights) <= 0 {
		return "", ""
	}

	// two characters for each rune, plus one for end of line
	chs := make([]byte, 0, 2*lastCh+1)

	i, h := 0, 0
	for i < lastCh && h < len(c.Highlights) && i <= c.Content.Length() {
		ch, _ := c.Content.Char(i)
		if i < c.Highlights[h].Start {
			switch {
			case ch == '\t':
				chs = append(chs, '\t')

			case unicode.In(ch, unicode.Han):
				chs = append(chs, ' ', ' ')

			default:
				chs = append(chs, ' ')
			}

			if !hasHighlighted {
				lead = len(chs)
			}
			i++

		} else {
			switch {
			case unicode.In(ch, unicode.Han):
				chs = append(chs, '^', '^')

			default:
				chs = append(chs, '^')
			}

			hasHighlighted = true
			i++
		}

		if i >= c.Highlights[h].End {
			h++
		}
	}

	// messageContent := fmt.Sprintf("%d:%d %s", c.Line+1, lead, message)
	return placeholder + string(chs), placeholder + string(chs[:lead]) + message
}

func (c *LineContext) Highlight(message string) string {
	lineNoWidth := 0
	lineNoStart := c.Content.Line
	lineNoFormat := ""
	lineNoPlaceholder := ""
	hasLineNo := c.Content.Line >= 0
	if hasLineNo {
		lineNoWidth = 1 + int(math.Log10(float64(lineNoStart)))
		lineNoFormat = fmt.Sprintf("%%%dd    ", lineNoWidth)
		lineNoPlaceholder = strings.Repeat(" ", lineNoWidth+4)
	}

	result := make([]string, 0, 3)
	if hasLineNo {
		line := c.MakeLineWithNumber(lineNoFormat)
		result = append(result, line)
		highlight, indentedMessage := c.MakeHighlight(lineNoPlaceholder, message)
		if len(highlight) > 0 {
			result = append(result, highlight, indentedMessage)
		}
	}

	return JoinLines(result)
}

func (c *LineContext) MakeContent() string {
	return c.Content.AsString()
}

type Context struct {
	Filename  string
	PrevLines []*LineContext
	NextLines []*LineContext
	Lines     []*LineContext
}

func NewContext(filename string, lines ...*LineContext) *Context {
	ctx := &Context{
		Filename: filename,
		Lines:    lines,
	}

	return ctx
}

func (c *Context) AddPrevLines(lines []*LineContext) {
	c.PrevLines = append(c.PrevLines, lines...)
}

func (c *Context) AddNextLines(lines []*LineContext) {
	c.NextLines = append(c.NextLines, lines...)
}

func (c *Context) highlight(format string, args ...any) []string {
	if len(c.Lines) <= 0 {
		return make([]string, 0)
	}

	message := fmt.Sprintf(format, args...)
	lastLine := c.Lines[len(c.Lines)-1]
	lineNoWidth := 0
	lineNoEnd := lastLine.Content.Line + len(c.NextLines) + 1
	lineNoFormat := ""
	lineNoPlaceholder := ""
	hasLineNo := lastLine.Content.Line >= 0
	if hasLineNo {
		lineNoWidth = 1 + int(math.Log10(float64(lineNoEnd)))
		lineNoFormat = fmt.Sprintf("%%%dd    ", lineNoWidth)
		lineNoPlaceholder = strings.Repeat(" ", lineNoWidth+4)
	}

	resultLines := len(c.PrevLines) + len(c.NextLines) + (3 * len(c.Lines))
	result := make([]string, 0, resultLines)

	for _, line := range c.PrevLines {
		content := line.MakeLineWithNumber(lineNoFormat)
		result = append(result, content)
	}

	for i, line := range c.Lines {
		lastLine := i == len(c.Lines)-1
		content := line.MakeLineWithNumber(lineNoFormat)
		result = append(result, content)
		highlight, indentedMessage := line.MakeHighlight(lineNoPlaceholder, message)
		if len(highlight) > 0 {
			result = append(result, highlight)
		}

		if lastLine && len(indentedMessage) > 0 {
			result = append(result, indentedMessage)
		}
	}

	return result
}

func (c *Context) Highlight(format string, args ...any) string {
	result := c.highlight(format, args...)
	return JoinLines(result)
}

func (c *Context) FullHighlight(format string, args ...any) string {
	result := c.highlight(format, args...)
	message := fmt.Sprintf("in '%s'", c.Filename)
	result = append(result, message)
	return JoinLines(result)
}

func JoinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

func JoinContexts(ctxList ...*Context) *Context {
	if len(ctxList) <= 0 {
		return nil
	}

	first := ctxList[0] // first context CAN BE nil
	result := &Context{}
	if first != nil {
		result.Filename = first.Filename
		result.PrevLines = make([]*LineContext, 0, len(first.PrevLines))
		result.PrevLines = append(result.PrevLines, first.PrevLines...)
		result.Lines = make([]*LineContext, 0, len(first.Lines)+len(ctxList)-1)
		result.Lines = append(result.Lines, first.Lines...)
	}

	lineMap := make(map[int]int)
	for i, line := range result.Lines {
		lineMap[line.Content.Line] = i
	}

	for i := 1; i < len(ctxList); i++ {
		ctx := ctxList[i]
		if ctx == nil {
			continue
		}

		if result.Filename != ctx.Filename {
			if len(result.Filename) <= 0 {
				result.Filename = ctx.Filename
			} else {
				return nil
			}
		}

		for _, line := range ctx.Lines {
			if lineIndex, ok := lineMap[line.Content.Line]; ok {
				for _, highlight := range line.Highlights {
					result.Lines[lineIndex].AddHighlight(highlight.Start, highlight.End)
					sort.Sort(ByHighlight(result.Lines[lineIndex].Highlights))
				}
			} else {
				result.Lines = append(result.Lines, line)
				lineMap[line.Content.Line] = len(result.Lines) - 1
			}
		}
	}

	return result
}
