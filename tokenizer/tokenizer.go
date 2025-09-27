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

func (t *Tokenizer) ScanWord(i int) (ast.Node, error) {
	r, _, eof := t.cursor.Peek(i)
	if eof || !IsValidIdentifierInitialRune(r) {
		return nil, nil
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
		return ast.NewIdentifier(ctx, content), nil
	}

	return ast.NewTerminalToken(ctx, tokenType), nil
}

func (t *Tokenizer) ScanFixedString(s string) *context.Context {
	return t.cursor.NextString(s)
}

func (t *Tokenizer) ScanSymbol() (ast.Node, error) {
	for _, op := range ast.OperatorList {
		ctx := t.cursor.NextString(op)
		if ctx != nil {
			tokenType := ast.GetOperatorTokenType(op)
			return ast.NewTerminalToken(ctx, tokenType), nil
		}
	}

	return nil, nil
}

func (t *Tokenizer) scanHexadecimalNumber() (ast.Node, error) {
	i := 2 // skip "0x"
	begin := t.cursor.State()
	v := uint64(0)
	for {
		r, eol, eof := t.cursor.Peek(i)
		if eol || eof {
			if i > 2 {
				break
			}

			// "0x" only
			state := t.cursor.PeekState(i)
			s, ctx := t.cursor.FinishWith(begin, state)
			return nil, ast.NewError(ctx, "invalid hexadecimal number '%s'", s)
		}

		if '0' <= r && r <= '9' {
			v = (v << 4) | uint64(r-'0')

		} else if 'a' <= r && r <= 'f' {
			v = (v << 4) | uint64(r-'a'+10)

		} else if 'A' <= r && r <= 'F' {
			v = (v << 4) | uint64(r-'A'+10)

		} else if ('g' <= r && r <= 'z') || ('G' <= r && r <= 'Z') {
			state := t.cursor.PeekState(i)
			s, ctx := t.cursor.FinishWith(begin, state)
			return nil, ast.NewError(ctx, "invalid hexadecimal number '%s'", s)

		} else {
			break
		}

		i++
	}

	if i > 18 {
		state := t.cursor.PeekState(i)
		s, ctx := t.cursor.FinishWith(begin, state)
		return nil, ast.NewError(ctx, "hexadecimal number '%s' is too large", s)
	}

	state := t.cursor.PeekState(i)
	_, ctx := t.cursor.FinishWith(begin, state)
	t.cursor.SetState(state)
	return ast.NewIntegerLiteral(ctx, v), nil
}

func (t *Tokenizer) ScanNumber() (ast.Node, error) {
	if ctx := t.cursor.PeekString("0x"); ctx != nil {
		return t.scanHexadecimalNumber()
	}

	return nil, nil
}

func (t *Tokenizer) ScanToken() (ast.Node, error) {
	r, _, eof := t.cursor.Rune()
	if eof {
		return nil, nil
	}

	if IsValidNumberInitialRune(r) {
		return t.ScanNumber()
	}

	if IsValidIdentifierInitialRune(r) {
		return t.ScanWord(0)
	}

	if IsValidSymbolRune(r) {
		return t.ScanSymbol()
	}

	return nil, nil
}
