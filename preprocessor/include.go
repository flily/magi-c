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
	content, contentCtx := cursorScanUntil(p.cursor, pos)
	if len(content) <= 0 {
		return nil, ast.NewError(contentCtx, "expected file name after '#include', got empty string")
	}

	rb, rbCtx := p.cursor.CurrentChar()
	if rb != pos {
		return nil, ast.NewError(rbCtx, "expected '%c' after file name in '#include', got '%c'", pos, rb)
	}

	directive.LBracket = lbCtx
	directive.Content = contentCtx
	directive.RBracket = rbCtx

	return directive, nil
}
