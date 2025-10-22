package preprocessor

import (
	"slices"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
)

type Preprocessor interface {
	Process(hash *context.Context, name *context.Context) (ast.Node, error)
}

type PreprocessorInitializer func(cursor *context.Cursor) Preprocessor

func cursorScanUntil(cursor *context.Cursor, flags ...rune) (string, *context.Context) {
	begin := cursor.State()
	for {
		r, eol, eof := cursor.Rune()
		if eol || eof {
			break
		}

		if slices.Contains(flags, r) {
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

func ScanDirective(cursor *context.Cursor) (string, *context.Context, *context.Context, error) {
	hash, hashCtx := cursor.CurrentChar()
	if hash != '#' {
		return "", nil, nil, ast.NewError(hashCtx, "expected '#' at the beginning of preprocessor directive, got '%c'", hash)
	}

	if !cursor.IsFirstNonWhiteChar() {
		return "", nil, nil, ast.NewError(hashCtx, "'#' must be the first non-whitespace character in the line")
	}

	cursor.NextInLine()
	name, nameCtx := cursorScanUntil(cursor, ' ', '\t')
	if len(name) <= 0 {
		return "", nil, nil, ast.NewError(nameCtx, "expected preprocessor directive name after '#', got empty string")
	}

	return name, hashCtx, nameCtx, nil
}

func scanDirectiveOn(cursor *context.Cursor, prep PreprocessorInitializer) (ast.Node, error) {
	_, ctxHash, ctxCmd, err := ScanDirective(cursor)
	if err != nil {
		return nil, err
	}

	p := prep(cursor)
	return p.Process(ctxHash, ctxCmd)
}
