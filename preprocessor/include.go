package preprocessor

import (
	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
)

type preprocessorInclude struct {
	cursor *context.Cursor
}

func Include(cursor *context.Cursor) Preprocessor {
	p := &preprocessorInclude{
		cursor: cursor,
	}

	return p
}

func (p *preprocessorInclude) Process(hash *context.Context, name *context.Context) (ast.Node, error) {
	directive := &ast.PreprocessorInclude{
		Hash:    hash,
		Command: name,
	}

	p.cursor.SkipWhitespaceInLine()
	lb, lbCtx := p.cursor.CurrentChar()
	pos := rune(0)
	switch lb {
	case '<':
		pos = '>'
	case '"':
		pos = '"'
	default:
		return nil, ast.NewError(lbCtx, "expected '<' or '\"' after '#include', got '%c'", lb)
	}

	p.cursor.NextInLine()
	content, contentCtx := cursorScanUntilInLine(p.cursor, '>', '"')
	rb, rbCtx := p.cursor.CurrentChar()
	if rb == 0 {
		ctx := context.Join(lbCtx, rbCtx)
		return nil, ast.NewError(ctx, "quote not closed")
	}

	if rb != pos {
		ctx := context.Join(lbCtx, rbCtx)
		return nil, ast.NewError(ctx, "quote mismatch, expected '%c', got '%c'", pos, rb)
	}

	if len(content) <= 0 {
		return nil, ast.NewError(lbCtx, "expected file name after '#include', got empty string")
	}

	directive.LBracket = lbCtx
	directive.Content = contentCtx
	directive.RBracket = rbCtx

	return directive, nil
}
