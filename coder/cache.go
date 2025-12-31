package coder

import (
	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/parser"
	"github.com/flily/magi-c/preprocessor"
	"github.com/flily/magi-c/tokenizer"
)

type Cache struct {
	Documents map[string]*ast.Document
}

func NewCache() *Cache {
	c := &Cache{
		Documents: make(map[string]*ast.Document),
	}

	return c
}

func (c *Cache) Add(index string, doc *ast.Document) {
	c.Documents[index] = doc
}

func (c *Cache) Get(index string) (*ast.Document, bool) {
	doc, ok := c.Documents[index]
	return doc, ok
}

func (c *Cache) NewCoderFromDocument(doc *ast.Document) *Coder {
	return NewCoder(c, doc)
}

func (c *Cache) NewCoderFromBinary(data []byte, filename string) (*Coder, error) {
	t := tokenizer.NewTokenizerFrom(data, filename)
	parser := parser.NewLLParser(t)
	preprocessor.RegisterPreprocessors(parser)
	doc, err := parser.Parse()
	if err != nil {
		return nil, err
	}

	return c.NewCoderFromDocument(doc), nil
}
