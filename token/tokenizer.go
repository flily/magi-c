package token

import (
	"os"
)

type TokenizerState int

const (
	TokenizerStateInit TokenizerState = iota
)

type TokenContext struct {
	Token   Token
	Context *LineContext
}

type Tokenizer struct {
	Filename string
	Lines    []*LineContext
	state    TokenizerState
	line     int
	column   int
}

func NewTokenizer() *Tokenizer {
	t := &Tokenizer{
		state: TokenizerStateInit,
		Lines: make([]*LineContext, 0),
	}

	return t
}

func (t *Tokenizer) ReadBuffer(buffer []byte, filename string) error {
	start := 0
	for i := 0; i < len(buffer); i++ {
		n := 0
		eolLF := buffer[i] == '\n'
		eolCRLF := buffer[i] == '\r' && i+1 < len(buffer) && buffer[i+1] == '\n'

		if eolLF || eolCRLF {
			line := NewLineFromBytes(len(t.Lines), buffer[start:i])
			lctx := NewLineContext(filename, line)
			t.Lines = append(t.Lines, lctx)

			if eolCRLF {
				n = 1
			}
		}

		i += n
	}
	return nil
}

func (t *Tokenizer) ReadFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return t.ReadBuffer(data, filename)
}

func (t *Tokenizer) nextRuneInLine() (rune, bool) {
	if t.line >= len(t.Lines) {
		return 0, true // No more lines
	}

	line := t.Lines[t.line]
	if t.column >= line.Content.Length() {
		return 0, true // No more characters in the current line
	}

	r := line.Content.Content[t.column]
	t.column++

	return r, false
}

func (t *Tokenizer) nextRune() (rune, bool) {
	if t.line >= len(t.Lines) {
		return 0, true // No more lines
	}

	line := t.Lines[t.line]
	r := line.Content.Content[t.column]

	t.column++
	for t.column >= line.Content.Length() {
		t.line++
		t.column = 0
		if t.line >= len(t.Lines) {
			return 0, true // No more lines
		}
		line = t.Lines[t.line]
	}

	return r, false
}

func (t *Tokenizer) currentRune() rune {
	line := t.Lines[t.line]
	return line.Content.Content[t.column]
}

func (t *Tokenizer) SkipWhitespace() {
}

func IsValidIdentifierRune(r rune) bool {
	if 'a' <= r && r <= 'z' {
		return true
	}

	if 'A' <= r && r <= 'Z' {
		return true
	}

	if '0' <= r && r <= '9' {
		return true
	}

	if r == '_' {
		return true
	}

	return false
}

func IsValidIdentifierInitialRune(r rune) bool {
	if 'a' <= r && r <= 'z' {
		return true
	}

	if 'A' <= r && r <= 'Z' {
		return true
	}

	if r == '_' {
		return true
	}

	return false
}

func (t *Tokenizer) scanIdentifier() *LineContext {
	line := t.Lines[t.line]
	start := t.column
	end := start
	for end < line.Content.Length() {
		c := line.Content.Rune(end)
		if !IsValidIdentifierRune(c) {
			break
		}

		end++
	}

	ctx := line.MarkContext(start, end)
	return ctx
}

func (t *Tokenizer) scanTokenInit(r rune) *TokenContext {
	if IsValidIdentifierInitialRune(r) {
		t.scanIdentifier()
	}

	return nil
}

func (t *Tokenizer) ScanToken() *TokenContext {
	t.SkipWhitespace()

	r := t.currentRune()
	switch t.state {
	case TokenizerStateInit:
		return t.scanTokenInit(r)

	default:
		return nil
	}
}
