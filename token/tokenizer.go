package token

import (
	"os"

	"github.com/flily/magi-c/context"
)

type TokenizerState int

const (
	TokenizerStateInit TokenizerState = iota
)

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
