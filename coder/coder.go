package coder

import (
	"github.com/flily/magi-c/ast"
)

type Coder struct {
	Refs     *Cache
	Document *ast.Document
}

func NewCoder(refs *Cache, doc *ast.Document) *Coder {
	c := &Coder{
		Refs:     refs,
		Document: doc,
	}

	return c
}
