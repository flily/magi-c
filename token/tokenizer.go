package token

import "os"

type TokenizerState int

const (
	TokenizerStateInit TokenizerState = iota
)

type Tokenizer struct {
	Filename string
	state    TokenizerState
	Lines    []*LineContext
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
