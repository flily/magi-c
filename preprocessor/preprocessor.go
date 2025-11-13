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

var preprocessors = map[string]PreprocessorInitializer{
	"include": Include,
	"inline":  Inline,
}

type PreprocessorRegistry interface {
	RegisterPreprocessor(name string, handler PreprocessorInitializer)
}

func RegisterPreprocessors(registry PreprocessorRegistry) {
	for name, handler := range preprocessors {
		registry.RegisterPreprocessor(name, handler)
	}
}

type cursorContainer struct {
	cursor *context.Cursor
}

func newCursorContainer(cursor *context.Cursor) cursorContainer {
	return cursorContainer{
		cursor: cursor,
	}
}

func cursorScanUntilInLine(cursor *context.Cursor, flags ...rune) (string, *context.Context) {
	if eol, eof := cursor.End(); eol || eof {
		_, ctx := cursor.CurrentChar()
		return "", ctx
	}

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

func isValidDirectiveNameChar(r rune) bool {
	if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9') {
		return true
	}

	if r == '_' || r == '-' {
		return true
	}

	return false
}

func scanDirectiveName(cursor *context.Cursor) (string, *context.Context) {
	begin := cursor.State()
	for {
		r, eol, eof := cursor.Rune()
		if eol || eof {
			break
		}

		if !isValidDirectiveNameChar(r) {
			break
		}

		cursor.NextInLine()
	}

	result, ctx := cursor.Finish(begin)
	return result, ctx
}

func ScanDirective(cursor *context.Cursor) (string, *context.Context, *context.Context, error) {
	hash, hashCtx := cursor.CurrentChar()
	if hash != '#' {
		return "", nil, nil, ast.NewError(hashCtx, "expect '#' at the beginning of preprocessor directive, got '%c'", hash)
	}

	if !cursor.IsFirstNonWhiteChar() {
		return "", nil, nil, ast.NewError(hashCtx, "'#' must be the first non-whitespace character in the line")
	}

	cursor.NextInLine()
	name, nameCtx := scanDirectiveName(cursor)
	if len(name) <= 0 {
		_, ctx := cursor.CurrentChar()
		return "", nil, nil, ast.NewError(ctx, "expect preprocessor directive name after '#', got empty string")
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
