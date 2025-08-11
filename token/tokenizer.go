package token

import (
	"os"

	"github.com/flily/magi-c/context"
)

type TokenizerState int

const (
	TokenizerStateInit TokenizerState = iota
)

type TokenContext struct {
	Token   Token
	Context *context.Context
}

type Tokenizer struct {
	Filename string
	file     *context.FileContext
	state    TokenizerState
	cursor   *Cursor
}

func NewTokenizer() *Tokenizer {
	t := &Tokenizer{
		state: TokenizerStateInit,
	}

	return t
}

func (t *Tokenizer) ReadBuffer(buffer []byte, filename string) {
	file := context.ReadFileData(filename, buffer)
	t.file = file
	t.Filename = filename
	t.cursor = NewCursor(file)
}

func (t *Tokenizer) ReadFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	t.ReadBuffer(data, filename)
	return nil
}

func (t *Tokenizer) Cursor() *Cursor {
	if t.cursor == nil {
		t.cursor = NewCursor(t.file)
	}

	return t.cursor
}
