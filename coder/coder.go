package coder

import (
	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/parser"
	"github.com/flily/magi-c/preprocessor"
	"github.com/flily/magi-c/tokenizer"
)

type Coder struct {
	Document *ast.Document
}

func NewCoderFrom(document *ast.Document) *Coder {
	c := &Coder{
		Document: document,
	}

	return c
}

func NewCoderFromBinary(data []byte, filename string) (*Coder, error) {
	t := tokenizer.NewTokenizerFrom(data, filename)
	parser := parser.NewLLParser(t)
	preprocessor.RegisterPreprocessors(parser)
	doc, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return NewCoderFrom(doc), nil
}
