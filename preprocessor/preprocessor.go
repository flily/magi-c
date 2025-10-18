package preprocessor

import (
	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
)

type Preprocessor interface {
	Process(hash *context.Context, name *context.Context) (ast.Node, error)
}

type PreprocessorInitializer func(cursor *context.Cursor) Preprocessor
