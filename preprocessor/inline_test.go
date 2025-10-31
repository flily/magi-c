package preprocessor

import (
	"strings"
	"testing"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
)

func TestInlineDirectiveBasic(t *testing.T) {
	code := strings.Join([]string{
		"#inline c",
		"#include <stdio.h>",
		"#inline asm",
		"#end-inline c",
	}, "\n")

	cursor := context.NewCursorFromString("example.c", code)
	node, err := scanDirectiveOn(cursor, Inline)
	if err != nil {
		t.Fatalf("unexpected error:\n%v", err)
	}

	result, ok := node.(*ast.PreprocessorInline)
	if !ok {
		t.Fatalf("expected PreprocessorInline node, got %T", node)
	}

	if result == nil {
		t.Fatalf("expected non-nil PreprocessorInline node")
	}

	gotContent := result.Content.HighlightText("here")
	expContent := strings.Join([]string{
		"   2:   #include <stdio.h>",
		"        ^^^^^^^^^^^^^^^^^^",
		"   3:   #inline asm",
		"        ^^^^^^^^^^^",
		"        here",
	}, "\n")
	if gotContent != expContent {
		t.Errorf("expected content context highlight:\n%s\ngot:\n%s", expContent, gotContent)
	}
}
