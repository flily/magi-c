package coder

import (
	"github.com/flily/magi-c/ast"
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
