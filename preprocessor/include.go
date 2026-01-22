package preprocessor

import (
	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
)

type preprocessorInclude struct {
	cursorContainer
}

func Include(cursor *context.Cursor) Preprocessor {
	p := &preprocessorInclude{
		cursorContainer: newCursorContainer(cursor),
	}

	return p
}

func (p *preprocessorInclude) Process(hash *context.Context, name *context.Context) (ast.TerminalNode, error) {
	p.cursor.SkipWhitespaceInLine()
	lb, lbCtx := p.cursor.CurrentChar()
	pos := rune(0)
	switch lb {
	case '<':
		pos = '>'
	case '"':
		pos = '"'
	default:
		return nil, lbCtx.Error("expected '<' or '\"' after '#include', got '%c'", lb)
	}

	p.cursor.NextInLine()
	content, contentCtx := cursorScanUntilInLine(p.cursor, '>', '"')
	rb, rbCtx := p.cursor.CurrentChar()
	if rb == 0 {
		ctx := context.Join(lbCtx, rbCtx)
		return nil, ctx.Error("quote not closed")
	}
	p.cursor.SkipInLine(1)

	if rb != pos {
		ctx := context.Join(lbCtx, rbCtx)
		return nil, ctx.Error("quote mismatch, expected '%c', got '%c'", pos, rb)
	}

	if len(content) <= 0 {
		return nil, lbCtx.Error("expected file name after '#include', got empty string")
	}

	directive := ast.NewPreprocessorInclude(hash, name, lbCtx, contentCtx, rbCtx)

	return directive, nil
}
