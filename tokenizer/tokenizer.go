package tokenizer

import (
	"os"

	"github.com/flily/magi-c/ast"
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

func (t *Tokenizer) ScanWord(i int) ast.Node {
	r, _, eof := t.cursor.Peek(i)
	if eof || !IsValidIdentifierInitialRune(r) {
		return nil
	}

	start := t.cursor.State()
	for IsValidIdentifierRune(r) && !eof {
		r, eol, eof := t.cursor.Peek(i)
		if eol || eof {
			break
		}

		if !IsValidIdentifierRune(r) {
			break
		}

		i++
	}

	finish := t.cursor.PeekState(i)
	t.cursor.SetState(finish)

	content, ctx := t.cursor.FinishWith(start, finish)
	tokenType := ast.GetKeywordTokenType(content)

	if tokenType == ast.Invalid {
		return ast.NewIdentifier(ctx, content)
	}

	return ast.NewTerminalToken(ctx, tokenType)
}

func (t *Tokenizer) ScanFixedString(s string) *context.Context {
	return t.cursor.NextString(s)
}

func (t *Tokenizer) ScanSymbol() ast.Node {
	for _, op := range ast.OperatorList {
		ctx := t.cursor.NextString(op)
		if ctx != nil {
			tokenType := ast.GetOperatorTokenType(op)
			return ast.NewTerminalToken(ctx, tokenType)
		}
	}

	return nil
}

func (t *Tokenizer) ScanToken() ast.Node {
	r, _, eof := t.cursor.Rune()
	if eof {
		return nil
	}

	if IsValidIdentifierInitialRune(r) {
		return t.ScanWord(0)
	}

	if IsValidSymbolRune(r) {
		return t.ScanSymbol()
	}

	return nil
}
