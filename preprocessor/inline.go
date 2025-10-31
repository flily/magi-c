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
	p.cursor.SkipWhitespaceInLine()
	blockType, btCtx := cursorScanUntilInLine(p.cursor, ' ', '\t')
	if len(blockType) <= 0 {
		return nil, ast.NewError(btCtx, "expect block type")
	}

	p.cursor.SkipWhitespaceInLine()
	if eol, _ := p.cursor.End(); !eol {
		content, ctx := cursorScanUntilInLine(p.cursor)
		return nil, ast.NewError(ctx, "expected EOL after inline block type, got '%s'", content)
	}

	eof := p.cursor.NextLine()
	if eof {
		_, ctx := p.cursor.CurrentChar()
		return nil, ast.NewError(ctx, "expected inline block content, got EOF")
	}

	var contentCtx *context.Context

	for {
		endName, endHashCtx, endNameCtx, err := ScanDirective(p.cursor)
		if endName == "end-inline" && err == nil {
			p.cursor.SkipWhitespaceInLine()
			endBlockType, endBtCtx := cursorScanUntilInLine(p.cursor, ' ', '\t')
			if endBlockType == blockType {
				directive := ast.NewPreprocessorInline(hash, name, btCtx, contentCtx, endHashCtx, endNameCtx, endBtCtx)
				return directive, nil
			}
		}

		_, lineCtx := p.cursor.CurrentLine()
		if contentCtx == nil {
			contentCtx = lineCtx
		} else {
			contentCtx = context.Join(contentCtx, lineCtx)
		}

		cursorScanUntilInLine(p.cursor)
		eof := p.cursor.NextLine()
		if eof {
			_, ctx := p.cursor.CurrentChar()
			return nil, ast.NewError(ctx, "expected '#inline %s' to close inline block, got EOF", blockType)
		}
	}
}
