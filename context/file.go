package context

import (
	"fmt"
	"os"
)

var (
	EolCR   = []byte{'\r'}
	EolLF   = []byte{'\n'}
	EolCRLF = []byte{'\r', '\n'}
)

func BytesToRunes(bs []byte) []rune {
	rs := make([]rune, len(bs))
	for i, b := range bs {
		rs[i] = rune(b)
	}
	return rs
}

func RunesToBytes(rs []rune) []byte {
	bs := make([]byte, len(rs))
	for i, r := range rs {
		bs[i] = byte(r)
	}
	return bs
}

type LineContent struct {
	Line    int
	EOL     []rune
	Content []rune
}

func NewLineFromBytes(line int, data []byte, eol []byte) LineContent {
	l := LineContent{
		Line:    line,
		EOL:     BytesToRunes(eol),
		Content: []rune(string(data)),
	}

	return l
}

func (l *LineContent) Rune(n int) rune {
	if n < 0 || n >= len(l.Content) {
		return 0
	}

	return l.Content[n]
}

func (l *LineContent) EOLBytes() []byte {
	return RunesToBytes(l.EOL)
}

func (l *LineContent) EOLString() string {
	if len(l.EOL) <= 0 {
		return "<EOF>"

	} else if len(l.EOL) == 1 && l.EOL[0] == '\n' {
		return "<EOL LF>"

	} else if len(l.EOL) == 1 && l.EOL[0] == '\r' {
		return "<EOL CR>"

	} else if len(l.EOL) == 2 && l.EOL[0] == '\r' && l.EOL[1] == '\n' {
		return "<EOL CRLF>"

	}

	parts := make([]string, 0, len(l.EOL))
	for _, r := range l.EOL {
		parts = append(parts, fmt.Sprintf("0x%X", r))
	}

	return "<EOL " + fmt.Sprintf("%s", parts) + ">"
}

func (l *LineContent) String() string {
	return string(l.Content)
}

func (l *LineContent) StringWithEOL() string {
	return l.String() + l.EOLString()
}

func (l *LineContent) Length() int {
	return len(l.Content)
}

func (l *LineContent) ToLineContext(file *FileContext) *LineContext {
	ctx := &LineContext{
		Content: l,
		File:    file,
	}

	return ctx
}

func (l *LineContent) MarkLine(start int, end int) *LineContext {
	if start > l.Length() || end > l.Length()+1 {
		err := fmt.Errorf("invalid context argument start=%d end=%d length=%d",
			start, end, l.Length())

		panic(err)
	}

	ctx := &LineContext{
		Content: l,
	}

	return ctx.MarkLine(start, end)
}

func (l *LineContent) Mark(start int, end int) (string, *Context) {
	lctx := l.MarkLine(start, end)
	content := ""
	if start < l.Length() {
		content = string(l.Content[start:end])
	}
	return content, lctx.ToContext()
}

type FileContext struct {
	Filename string
	Contents []*LineContent
}

// meetEOL checks if current position is End of Line. A non-negative value indicates EOL.
func meetEOL(data []byte, i int) (int, []byte) {
	if i >= len(data) {
		return 0, nil
	}

	if data[i] == '\n' {
		return 1, EolLF
	}

	if data[i] == '\r' {
		if i+1 >= len(data) {
			return 1, EolCR
		}

		if data[i+1] == '\n' {
			return 2, EolCRLF
		}
	}

	return -1, nil
}

func ReadFileData(filename string, data []byte) *FileContext {
	lines := make([]*LineContent, 0, 64)
	column := 0
	i := 0
	start := 0
	for i < len(data) {
		n, eol := meetEOL(data, i)
		if n < 0 {
			i++
			column++
			continue
		}

		lineContent := NewLineFromBytes(len(lines), data[start:i], eol)
		lines = append(lines, &lineContent)
		i += n
		column = 0
		start = i
	}

	if i > start {
		lineContent := NewLineFromBytes(len(lines), data[start:i], nil)
		lines = append(lines, &lineContent)
	}

	ctx := &FileContext{
		Filename: filename,
		Contents: lines,
	}

	return ctx
}

func ReadFile(filename string) (*FileContext, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return ReadFileData(filename, data), nil
}

func (f *FileContext) Rune(line int, column int) (rune, bool) {
	if line >= len(f.Contents) {
		return 0, true
	}

	l := f.Line(line)
	if column >= len(l.Content) {
		return 0, true
	}

	return l.Content[column], false
}

func (f *FileContext) Lines() int {
	return len(f.Contents)
}

func (f *FileContext) Line(n int) *LineContent {
	if n < 0 || n >= len(f.Contents) {
		return nil
	}

	return f.Contents[n]
}

func (f *FileContext) LineContext(n int) *LineContext {
	ctx := &LineContext{
		Content: f.Line(n),
		File:    f,
	}

	return ctx
}
