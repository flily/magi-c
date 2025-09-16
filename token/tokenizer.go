package token

import (
	"os"

	"github.com/flily/magi-c/context"
)

type TokenizerState int

const (
	TokenizerStateInit TokenizerState = iota
)

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
		_, _ = t.cursor.NextInLine()
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

func (t *Tokenizer) ScanToken() *context.Context {
	r, _, eof := t.cursor.Rune()
	if eof {
		return nil
	}

	if IsValidIdentifierInitialRune(r) {
		return t.ScanWord()
	}

	if IsValidSymbolRune(r) {
		return t.ScanSymbol()
	}

	return nil
}
