package preprocessor

import (
	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
)

type preprocessorInline struct {
	cursorContainer
}

func Inline(cursor *context.Cursor) Preprocessor {
	p := &preprocessorInline{
		cursorContainer: newCursorContainer(cursor),
	}

	return p
}

func (p *preprocessorInline) Process(hash *context.Context, name *context.Context) (ast.Node, error) {
	return nil, nil
}
