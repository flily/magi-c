package token

import (
	"os"

	"github.com/flily/magi-c/context"
)

type TokenizerState int

const (
	TokenizerStateInit TokenizerState = iota
)

func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
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

var validSymbolInitialRunes = []bool{
	'!':  true,
	'"':  true,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'(':  true,
	')':  true,
	'*':  true,
	'+':  true,
	',':  true,
	'-':  true,
	'.':  true,
	'/':  true,
	'[':  true,
	'\\': true,
	']':  true,
	'^':  true,
	'~':  true,
}

func IsValidSymbolRune(r rune) bool {
	if int(r) < len(validSymbolInitialRunes) {
		return validSymbolInitialRunes[int(r)]
	}

	return false
}

type Tokenizer struct {
	Filename string
	state    TokenizerState
	cursor   *Cursor
}

func NewTokenizerFrom(buffer []byte, filename string) *Tokenizer {
	file := context.ReadFileData(filename, buffer)

	cursor := NewCursor(file)
	t := &Tokenizer{
		Filename: filename,
		state:    TokenizerStateInit,
		cursor:   cursor,
	}

	return t
}

func NewTokenizerFromFile(filename string) (*Tokenizer, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return NewTokenizerFrom(file, filename), nil
}

func (t *Tokenizer) CurrentChar() *context.Context {
	return t.cursor.CurrentChar()
}

func (t *Tokenizer) SkipWhitespace() {
	for {
		_, eol, _ := t.cursor.Rune()
		if eol {
			if eof := t.cursor.NextNonEmptyLine(); eof {
				return
			}

			continue
		}

		r, _, _ := t.cursor.Rune()
		if IsWhitespace(r) {
			t.cursor.Next()

		} else {
			break
		}
	}
}

func (t *Tokenizer) ScanWord() *context.Context {
	r, _, eof := t.cursor.Rune()
	if eof || !IsValidIdentifierInitialRune(r) {
		return nil
	}

	start := t.cursor.State()
	for IsValidIdentifierRune(r) && !eof {
		t.cursor.Next()
		r, _, eof = t.cursor.Rune()
	}

	return t.cursor.Finish(start)
}

func (t *Tokenizer) ScanFixedString(s string) *context.Context {
	return t.cursor.NextString(s)
}

func (t *Tokenizer) ScanSymbol() *context.Context {
	for _, op := range operatorList {
		ctx := t.cursor.NextString(op)
		if ctx != nil {
			return ctx
		}
	}

	return nil
}
