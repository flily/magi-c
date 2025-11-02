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
		t.Fatalf("expect PreprocessorInline node, got %T", node)
	}

	if result == nil {
		t.Fatalf("expect non-nil PreprocessorInline node")
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
		t.Errorf("expect content context highlight:\n%s\ngot:\n%s", expContent, gotContent)
	}
}

func TestInlineDirectiveWithEmptyBlock(t *testing.T) {
	code := strings.Join([]string{
		"#inline c",
		"#end-inline c",
	}, "\n")

	cursor := context.NewCursorFromString("example.c", code)
	node, err := scanDirectiveOn(cursor, Inline)
	if err != nil {
		t.Fatalf("unexpected error:\n%v", err)
	}

	result, ok := node.(*ast.PreprocessorInline)
	if !ok {
		t.Fatalf("expect PreprocessorInline node, got %T", node)
	}

	if result == nil {
		t.Fatalf("expect non-nil PreprocessorInline node")
	}

	if result.Content != nil {
		t.Errorf("expect nil content context for empty inline block, got non-nil")
	}

	if !result.Empty() {
		t.Errorf("expect empty node returned")
	}
}

func TestInlineDirectiveWithoutBlockType(t *testing.T) {
	code := strings.Join([]string{
		"#inline",
		"#include <stdio.h>",
		"#end-inline",
	}, "\n")

	cursor := context.NewCursorFromString("example.c", code)
	_, err := scanDirectiveOn(cursor, Inline)

	got := err.Error()
	exp := strings.Join([]string{
		"   1:   #inline<EOL LF>",
		"               ^^^^^^^^",
		"               expect block type",
	}, "\n")

	if got != exp {
		t.Errorf("expected error message:\n%s\ngot:\n%s", exp, got)
	}
}

func TestInlineDirectiveWithWrongBlockType(t *testing.T) {
	code := strings.Join([]string{
		"#inline c asm",
		"#include <stdio.h>",
		"#end-inline c",
	}, "\n")

	cursor := context.NewCursorFromString("example.c", code)
	_, err := scanDirectiveOn(cursor, Inline)

	got := err.Error()
	exp := strings.Join([]string{
		"   1:   #inline c asm",
		"                  ^^^",
		"                  expected EOL after inline block type, got 'asm'",
	}, "\n")

	if got != exp {
		t.Errorf("expected error message:\n%s\ngot:\n%s", exp, got)
	}
}

func TestInlineDirectiveWithNoContentAndNoClosing(t *testing.T) {
	code := strings.Join([]string{
		"#inline c",
	}, "\n")

	cursor := context.NewCursorFromString("example.c", code)
	_, err := scanDirectiveOn(cursor, Inline)

	got := err.Error()
	exp := strings.Join([]string{
		"   1:   #inline c<EOF>",
		"                 ^^^^^",
		"                 expect inline block content, got EOF",
	}, "\n")

	if got != exp {
		t.Errorf("expected error message:\n%s\ngot:\n%s", exp, got)
	}
}

func TestInlineDirectiveWithUnclosedBlock(t *testing.T) {
	code := strings.Join([]string{
		"#inline c",
		"#include <stdio.h>",
	}, "\n")

	cursor := context.NewCursorFromString("example.c", code)
	_, err := scanDirectiveOn(cursor, Inline)

	got := err.Error()
	exp := strings.Join([]string{
		"   2:   #include <stdio.h><EOF>",
		"                          ^^^^^",
		"                          expect '#end-inline c' to close inline block, got EOF",
	}, "\n")

	if got != exp {
		t.Errorf("expected error message:\n%s\ngot:\n%s", exp, got)
	}
}
