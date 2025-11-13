package parser

import (
	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/preprocessor"
)

type Parser interface {
	Parse() (*ast.Document, error)
	RegisterPreprocessor(name string, handler preprocessor.PreprocessorInitializer)
}
