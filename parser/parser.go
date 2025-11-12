package parser

import (
	"github.com/flily/magi-c/ast"
)

type Parser interface {
	Parse() (*ast.Document, error)
}
