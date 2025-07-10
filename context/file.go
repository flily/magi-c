package context

import (
	"os"
)

type LineContent struct {
	Line    int
	Content []rune
}

func NewLineFromBytes(line int, data []byte) LineContent {
	l := LineContent{
		Line:    line,
		Content: []rune(string(data)),
	}

	return l
}

func (l LineContent) String() string {
	return string(l.Content)
}

type FileContext struct {
	Filename string
	Contents []*LineContent
}

// meetEOL checks if current position is End of Line. A non-negative value indicates EOL.
func meetEOL(data []byte, i int) int {
	if i >= len(data) {
		return 0
	}

	if data[i] == '\n' {
		return 1
	}

	if data[i] == '\r' {
		if i+1 >= len(data) {
			return 1
		}

		if data[i+1] == '\n' {
			return 2
		}
	}

	return -1
}

func ReadFileData(filename string, data []byte) *FileContext {
	lines := make([]*LineContent, 0, 64)
	line := 0
	column := 0
	i := 0
	start := 0
	for i < len(data) {
		n := meetEOL(data, i)
		if n < 0 {
			i++
			column++
			continue
		}

		lineContent := NewLineFromBytes(line, data[start:i])
		lines = append(lines, &lineContent)
		i += n
		column = 0
		start = i
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

func (f *FileContext) Lines() int {
	return len(f.Contents)
}

func (f *FileContext) LineContent(n int) *LineContent {
	if n < 0 || n >= len(f.Contents) {
		return nil
	}

	return f.Contents[n]
}

func (f *FileContext) LineContext(n int) *LineContext {
	ctx := &LineContext{
		Content: f.LineContent(n),
		File:    f,
	}

	return ctx
}
