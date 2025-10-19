package preprocessor

import (
	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
)

type Preprocessor interface {
	Process(hash *context.Context, name *context.Context) (ast.Node, error)
}

type PreprocessorInitializer func(cursor *context.Cursor) Preprocessor

func cursorScanUntil(cursor *context.Cursor, flag rune) (string, *context.Context) {
	begin := cursor.State()
	for {
		r, eol, eof := cursor.Rune()
		if eol || eof {
			break
		}

		if r == flag {
			break
		}

		cursor.NextInLine()
	}

	result, ctx := cursor.Finish(begin)
	return result, ctx
}

func (p *preprocessorInclude) Type() ast.NodeType {
	return ast.NodePreprocessorInclude
}
