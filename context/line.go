package context

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
)

const (
	FixedLeadingSpace  = "        "
	DefaultIndicator   = "^"
	NoHighlightMessage = ""
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

func (l *LineContext) Duplicate() *LineContext {
	dup := &LineContext{
		Content:    l.Content,
		File:       l.File,
		Highlights: make([]Highlight, len(l.Highlights)),
	}

	copy(dup.Highlights, l.Highlights)
	return dup
}

func (l *LineContext) ToContext() *Context {
	ctx := &Context{
		File: l.File,
		Lines: []*LineContext{
			l,
		},
	}

	return ctx
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

func (l *LineContext) IsSameLine(other *LineContext) bool {
	if l.File != other.File {
		return false
	}

	if l.Content.Line != other.Content.Line {
		return false
	}

	return true
}

func (l *LineContext) MarkLine(start int, end int) *LineContext {
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

func (l *LineContext) Mark(start int, end int) *Context {
	l.MarkLine(start, end)
	ctx := &Context{
		File:  l.File,
		Lines: []*LineContext{l},
	}

	return ctx
}

func (l *LineContext) LineNumber() string {
	return fmt.Sprintf("%4d:   ", l.Content.Line+1)
}

func (l *LineContext) String() string {
	return l.LineNumber() + l.Content.String()
}

func repeatToLength(s string, length int) string {
	l := len(s)
	if l >= length {
		return s[:length]
	}

	repeatCount := length / l
	remainder := length % l
	return strings.Repeat(s, repeatCount) + s[:remainder]
}

func (l *LineContext) HighlighTextWith(indicator string, format string, args ...any) string {
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
			repeatToLength(indicator, widthHighligh),
		)
	}

	content := l.String()
	highlight := fmt.Sprintf("%s\n%s%s",
		content,
		FixedLeadingSpace, strings.Join(parts, ""),
	)

	message := ""

	if len(format) > 0 {
		message += "\n" +
			FixedLeadingSpace +
			strings.Repeat(" ", len(lead)) +
			fmt.Sprintf(format, args...)
	}

	return highlight + message
}

func (l *LineContext) HighlighText(format string, args ...any) string {
	return l.HighlighTextWith(DefaultIndicator, format, args...)
}

func (l *LineContext) HighlightColour(colour color.Color, format string, args ...any) string {
	message := fmt.Sprintf(format, args...)

	parts := make([]string, 0, 8+2*len(l.Highlights))
	last, lead := 0, ""
	parts = append(parts, FixedLeadingSpace, l.StringContent(), "\n")

	for i, highlight := range l.Highlights {
		// highlight will store in order
		if highlight.Start < 0 || highlight.End > l.Length() || highlight.Start > highlight.End {
			err := fmt.Errorf("invalid highlight range: start=%d, end=%d, length=%d", highlight.Start, highlight.End, l.Length())
			panic(err)
		}

		parts = append(parts,
			string(l.Content.Content[last:highlight.Start]),
			colour.Sprint(string(l.Content.Content[highlight.Start:highlight.End])),
		)

		if i == 0 {
			widthSpace := 0
			for j := last; j < highlight.Start; j++ {
				widthSpace += CharWidthIn(l.Content.Content[j], j)
			}
			lead = strings.Repeat(" ", widthSpace)
		}

		last = highlight.End
	}

	return fmt.Sprintf("%s%s\n%s%s%s",
		FixedLeadingSpace, strings.Join(parts, ""),
		FixedLeadingSpace, strings.Repeat(" ", len(lead)), message,
	)
}

func (l *LineContext) Last() int {
	last := -1
	for _, h := range l.Highlights {
		if h.End > last {
			last = h.End
		}
	}

	return last
}

func FindLineContextSameLine(list []*LineContext, lctx *LineContext) *LineContext {
	for _, line := range list {
		if line.IsSameLine(lctx) {
			return line
		}
	}

	return nil
}
