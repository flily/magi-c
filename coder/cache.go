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

func (c *Cache) ParseContent(data []byte, filename string) (*ast.Document, error) {
	t := tokenizer.NewTokenizerFrom(data, filename)
	parser := parser.NewLLParser(t)
	preprocessor.RegisterPreprocessors(parser)
	return parser.Parse()
}
