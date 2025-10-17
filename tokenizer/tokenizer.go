package tokenizer

import (
	"math"
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
	cursor   *context.Cursor
}

func NewTokenizerFrom(buffer []byte, filename string) *Tokenizer {
	file := context.ReadFileData(filename, buffer)

	cursor := context.NewCursor(file)
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

func (t *Tokenizer) CurrentChar() (rune, *context.Context) {
	return t.cursor.CurrentChar()
}

func (t *Tokenizer) SkipWhitespaceInLine() {
	t.cursor.SkipWhitespaceInLine()
}

func (t *Tokenizer) SkipWhitespace() {
	t.cursor.SkipWhitespace()
}

func (t *Tokenizer) scanWord(i int) (string, *context.Context) {
	r, _, eof := t.cursor.Peek(i)
	if eof || !IsValidIdentifierInitialRune(r) {
		return "", nil
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
	return content, ctx
}

func (t *Tokenizer) ScanWordToken(i int) ast.Node {
	content, ctx := t.scanWord(i)

	tokenType := ast.GetKeywordTokenType(content)
	if tokenType == ast.Invalid {
		return ast.NewIdentifier(ctx, content)
	}

	return ast.NewTerminalToken(ctx, tokenType)
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

	_, ctx := t.cursor.CurrentChar()
	return nil, ast.NewError(ctx, "invalid symbol '%s'", ctx.Content())
}

func (t *Tokenizer) scanHexadecimalNumber() (ast.Node, error) {
	i := 2 // skip "0x"
	begin := t.cursor.State()
	v := uint64(0)
	invalidFormat := false

	for {
		r, eol, eof := t.cursor.Peek(i)
		if eol || eof {
			if i <= 2 {
				// "0x" only
				invalidFormat = true
			}

			break
		}

		if '0' <= r && r <= '9' {
			v = (v << 4) | uint64(r-'0')

		} else if 'a' <= r && r <= 'f' {
			v = (v << 4) | uint64(r-'a'+10)

		} else if 'A' <= r && r <= 'F' {
			v = (v << 4) | uint64(r-'A'+10)

		} else if ('g' <= r && r <= 'z') || ('G' <= r && r <= 'Z') {
			invalidFormat = true

		} else {
			break
		}

		i++
	}

	if invalidFormat {
		state := t.cursor.PeekState(i)
		s, ctx := t.cursor.FinishWith(begin, state)
		return nil, ast.NewError(ctx, "invalid hexadecimal number '%s'", s)
	}

	if i > 2+16 {
		state := t.cursor.PeekState(i)
		s, ctx := t.cursor.FinishWith(begin, state)
		return nil, ast.NewError(ctx, "hexadecimal number '%s' is too large", s)
	}

	state := t.cursor.PeekState(i)
	_, ctx := t.cursor.FinishWith(begin, state)
	t.cursor.SetState(state)
	return ast.NewIntegerLiteral(ctx, v), nil
}

func (t *Tokenizer) scanOctalNumber() (ast.Node, error) {
	i := 1 // skip "0"
	begin := t.cursor.State()
	v := uint64(0)
	invalidFormat := false

	for {
		r, eol, eof := t.cursor.Peek(i)
		if eol || eof {
			break
		}

		if '0' <= r && r <= '7' {
			v = (v << 3) | uint64(r-'0')

		} else if r == '8' || r == '9' || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
			invalidFormat = true

		} else {
			break
		}

		i++
	}

	if invalidFormat {
		state := t.cursor.PeekState(i)
		s, ctx := t.cursor.FinishWith(begin, state)
		return nil, ast.NewError(ctx, "invalid octal number '%s'", s)
	}

	if i > 1+21 {
		state := t.cursor.PeekState(i)
		s, ctx := t.cursor.FinishWith(begin, state)
		return nil, ast.NewError(ctx, "octal number '%s' is too large", s)
	}

	state := t.cursor.PeekState(i)
	_, ctx := t.cursor.FinishWith(begin, state)
	t.cursor.SetState(state)
	return ast.NewIntegerLiteral(ctx, v), nil
}

func (t *Tokenizer) scanDecimalNumber() (ast.Node, error) {
	i := 0
	begin := t.cursor.State()
	dotIndex, expIndex := -1, -1
	negExp := false
	integer, fraction, exponent := uint64(0), uint64(0), int(0)
	fractionExp := 1
	invalidFormat := false

	for {
		r, eol, eof := t.cursor.Peek(i)
		if eol || eof {
			// i MUST NOT be 0 here
			break
		}

		if dotIndex < 0 && expIndex < 0 {
			// integer part
			if '0' <= r && r <= '9' {
				integer = integer*10 + uint64(r-'0')

			} else if r == '.' {
				dotIndex = i

			} else if r == 'e' || r == 'E' {
				expIndex = i

			} else if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
				invalidFormat = true

			} else {
				break
			}

		} else if dotIndex >= 0 && expIndex < 0 {
			// fraction part
			if '0' <= r && r <= '9' {
				fraction = fraction*10 + uint64(r-'0')
				fractionExp *= 10

			} else if r == 'e' || r == 'E' {
				expIndex = i

			} else if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
				invalidFormat = true

			} else {
				break
			}

		} else if expIndex > 0 {
			// exponent part
			if '0' <= r && r <= '9' {
				exponent = (exponent * 10) + int(r-'0')

			} else if (i == expIndex+1) && (r == '+' || r == '-') {
				if r == '-' {
					negExp = true
				}

			} else if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
				invalidFormat = true

			} else {
				break
			}
		}

		i++
	}

	if invalidFormat {
		state := t.cursor.PeekState(i)
		s, ctx := t.cursor.FinishWith(begin, state)
		return nil, ast.NewError(ctx, "invalid decimal number '%s'", s)
	}

	state := t.cursor.PeekState(i)
	_, ctx := t.cursor.FinishWith(begin, state)
	t.cursor.SetState(state)

	if dotIndex < 0 && expIndex < 0 {
		// integer only
		return ast.NewIntegerLiteral(ctx, integer), nil

	} else if dotIndex >= 0 && expIndex < 0 {
		// fraction only
		fraction := float64(fraction) / float64(fractionExp)
		return ast.NewFloatLiteral(ctx, float64(integer)+fraction), nil

	} else {
		fraction := float64(fraction) / float64(fractionExp)
		base := float64(integer) + fraction

		if negExp {
			exponent = -exponent
		}

		return ast.NewFloatLiteral(ctx, math.Pow10(int(exponent))*base), nil
	}
}

func (t *Tokenizer) ScanNumber() (ast.Node, error) {
	r0, _, _ := t.cursor.Rune()
	begin := t.cursor.State()

	if r0 != '0' {
		return t.scanDecimalNumber()
	}

	r1, eol, eof := t.cursor.Peek(1)
	if eol || eof {
		_, ctx := t.cursor.Finish(begin)
		return ast.NewIntegerLiteral(ctx, 0), nil
	}

	if r1 == 'x' {
		return t.scanHexadecimalNumber()

	} else if '0' <= r1 && r1 <= '7' {
		return t.scanOctalNumber()

	} else {
		return t.scanDecimalNumber()
	}
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
		return t.ScanWordToken(0), nil
	}

	if IsValidSymbolRune(r) {
		return t.ScanSymbol()
	}

	return nil, nil
}
