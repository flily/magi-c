package preprocessor

import (
	"strings"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
)

const (
	PreprocessorCommandInline      = "inline"
	PreprocessorCommandInlineClose = "end-inline"
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

func (p *preprocessorInline) Process(hash *context.Context, name *context.Context) (ast.TerminalNode, error) {
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
		return nil, ast.NewError(ctx, "expect inline block content, got EOF")
	}

	var contentCtx *context.Context
	content := make([]string, 0, 64)

	for {
		endName, endHashCtx, endNameCtx, err := ScanDirective(p.cursor)
		if endName == PreprocessorCommandInlineClose && err == nil {
			p.cursor.SkipWhitespaceInLine()
			endBlockType, endBtCtx := cursorScanUntilInLine(p.cursor, ' ', '\t')
			if endBlockType == blockType {
				directive := ast.NewPreprocessorInline(hash, name, blockType, btCtx, strings.Join(content, "\n"), contentCtx, endHashCtx, endNameCtx, endBtCtx)
				return directive, nil
			}
		}

		_, lineCtx := p.cursor.CurrentLine()
		if contentCtx == nil {
			contentCtx = lineCtx
		} else {
			contentCtx = context.Join(contentCtx, lineCtx)
		}
		content = append(content, lineCtx.Content())

		cursorScanUntilInLine(p.cursor)
		eof := p.cursor.NextLine()
		if eof {
			_, ctx := p.cursor.CurrentChar()
			return nil, ast.NewError(ctx, "expect '#%s %s' to close inline block, got EOF", PreprocessorCommandInlineClose, blockType)
		}
	}
}
