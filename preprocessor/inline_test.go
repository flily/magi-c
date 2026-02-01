package preprocessor

import (
	"testing"

	"strings"

	"github.com/flily/magi-c/ast"
)

func TestInlineDirectiveBasic(t *testing.T) {
	code := strings.Join([]string{
		"#inline c",
		"#include <stdio.h>",
		"#inline asm",
		"#end-inline c",
	}, "\n")

	node, final := testScanDirectiveCorrect(t, code, Inline)
	result, ok := node.(*ast.PreprocessorInline)
	if !ok {
		t.Fatalf("expect PreprocessorInline node, got %T", node)
	}

	if result == nil {
		t.Fatalf("expect non-nil PreprocessorInline node")
		return // unreachable return, reduce lint error
	}

	expHash := strings.Join([]string{
		"    1 | #inline c",
		"      | ^",
		"      | here",
	}, "\n")
	checkElementContext(t, result.Hash, expHash)

	expCmd := strings.Join([]string{
		"    1 | #inline c",
		"      |  ^^^^^^",
		"      |  here",
	}, "\n")
	checkElementContext(t, result.Command, expCmd)

	expType := strings.Join([]string{
		"    1 | #inline c",
		"      |         ^",
		"      |         here",
	}, "\n")
	checkElementContext(t, result.CodeTypeCtx, expType)

	expContent := strings.Join([]string{
		"    2 | #include <stdio.h>",
		"      | ^^^^^^^^^^^^^^^^^^",
		"    3 | #inline asm",
		"      | ^^^^^^^^^^^",
		"      | here",
	}, "\n")
	checkElementContext(t, result.ContentCtx, expContent)

	expEndHash := strings.Join([]string{
		"    4 | #end-inline c",
		"      | ^",
		"      | here",
	}, "\n")
	checkElementContext(t, result.HashEnd, expEndHash)

	expEndCmd := strings.Join([]string{
		"    4 | #end-inline c",
		"      |  ^^^^^^^^^^",
		"      |  here",
	}, "\n")
	checkElementContext(t, result.CommandEnd, expEndCmd)

	expEndType := strings.Join([]string{
		"    4 | #end-inline c",
		"      |             ^",
		"      |             here",
	}, "\n")
	checkElementContext(t, result.CodeTypeEnd, expEndType)

	if result.Empty() {
		t.Errorf("expect non-empty node returned")
	}

	finalExp := strings.Join([]string{
		"    4 | #end-inline c<EOF>",
		"      |              ^^^^^",
		"      |              here",
	}, "\n")
	checkElementContext(t, final, finalExp)
}

func TestInlineDirectiveWithEmptyBlock(t *testing.T) {
	code := strings.Join([]string{
		"#inline c",
		"#end-inline c",
	}, "\n")

	node, final := testScanDirectiveCorrect(t, code, Inline)
	result, ok := node.(*ast.PreprocessorInline)
	if !ok {
		t.Fatalf("expect PreprocessorInline node, got %T", node)
	}

	expHash := strings.Join([]string{
		"    1 | #inline c",
		"      | ^",
		"      | here",
	}, "\n")
	checkElementContext(t, result.Hash, expHash)

	expCmd := strings.Join([]string{
		"    1 | #inline c",
		"      |  ^^^^^^",
		"      |  here",
	}, "\n")
	checkElementContext(t, result.Command, expCmd)

	expType := strings.Join([]string{
		"    1 | #inline c",
		"      |         ^",
		"      |         here",
	}, "\n")
	checkElementContext(t, result.CodeTypeCtx, expType)

	if result.ContentCtx != nil {
		t.Errorf("expect nil content context for empty inline block, got non-nil")
	}

	if !result.Empty() {
		t.Errorf("expect empty node returned")
	}

	expEndHash := strings.Join([]string{
		"    2 | #end-inline c",
		"      | ^",
		"      | here",
	}, "\n")
	checkElementContext(t, result.HashEnd, expEndHash)

	expEndCmd := strings.Join([]string{
		"    2 | #end-inline c",
		"      |  ^^^^^^^^^^",
		"      |  here",
	}, "\n")
	checkElementContext(t, result.CommandEnd, expEndCmd)

	expEndType := strings.Join([]string{
		"    2 | #end-inline c",
		"      |             ^",
		"      |             here",
	}, "\n")
	checkElementContext(t, result.CodeTypeEnd, expEndType)

	finalExp := strings.Join([]string{
		"    2 | #end-inline c<EOF>",
		"      |              ^^^^^",
		"      |              here",
	}, "\n")
	checkElementContext(t, final, finalExp)
}

func TestInlineDirectiveWithIndent(t *testing.T) {
	code := strings.Join([]string{
		"    #inline c",
		`    printf("hello, world\n");`,
		"    #end-inline c",
	}, "\n")

	node, final := testScanDirectiveCorrect(t, code, Inline)
	result, ok := node.(*ast.PreprocessorInline)
	if !ok {
		t.Fatalf("expect PreprocessorInline node, got %T", node)
	}

	expHash := strings.Join([]string{
		"    1 |     #inline c",
		"      |     ^",
		"      |     here",
	}, "\n")
	checkElementContext(t, result.Hash, expHash)

	expCmd := strings.Join([]string{
		"    1 |     #inline c",
		"      |      ^^^^^^",
		"      |      here",
	}, "\n")
	checkElementContext(t, result.Command, expCmd)

	expType := strings.Join([]string{
		"    1 |     #inline c",
		"      |             ^",
		"      |             here",
	}, "\n")
	checkElementContext(t, result.CodeTypeCtx, expType)

	expContent := strings.Join([]string{
		`    2 |     printf("hello, world\n");`,
		"      | ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^",
		"      | here",
	}, "\n")
	checkElementContext(t, result.ContentCtx, expContent)

	expEndHash := strings.Join([]string{
		"    3 |     #end-inline c",
		"      |     ^",
		"      |     here",
	}, "\n")
	checkElementContext(t, result.HashEnd, expEndHash)

	expEndCmd := strings.Join([]string{
		"    3 |     #end-inline c",
		"      |      ^^^^^^^^^^",
		"      |      here",
	}, "\n")
	checkElementContext(t, result.CommandEnd, expEndCmd)

	expEndType := strings.Join([]string{
		"    3 |     #end-inline c",
		"      |                 ^",
		"      |                 here",
	}, "\n")
	checkElementContext(t, result.CodeTypeEnd, expEndType)

	finalExp := strings.Join([]string{
		"    3 |     #end-inline c<EOF>",
		"      |                  ^^^^^",
		"      |                  here",
	}, "\n")
	checkElementContext(t, final, finalExp)
}

func TestInlineDirectiveWithoutBlockType(t *testing.T) {
	code := strings.Join([]string{
		"#inline",
		"#include <stdio.h>",
		"#end-inline",
	}, "\n")

	exp := strings.Join([]string{
		"example.mc:1:8: error: expect block type",
		"    1 | #inline<EOL LF>",
		"      |        ^^^^^^^^",
	}, "\n")
	checkScanDirectiveError(t, code, Inline, exp)
}

func TestInlineDirectiveWithWrongBlockType(t *testing.T) {
	code := strings.Join([]string{
		"#inline c asm",
		"#include <stdio.h>",
		"#end-inline c",
	}, "\n")

	exp := strings.Join([]string{
		"example.mc:1:11: error: expected EOL after inline block type, got 'asm'",
		"    1 | #inline c asm",
		"      |           ^^^",
	}, "\n")
	checkScanDirectiveError(t, code, Inline, exp)
}

func TestInlineDirectiveWithNoContentAndNoClosing(t *testing.T) {
	code := strings.Join([]string{
		"#inline c",
	}, "\n")

	exp := strings.Join([]string{
		"example.mc:1:10: error: expect inline block content, got EOF",
		"    1 | #inline c<EOF>",
		"      |          ^^^^^",
	}, "\n")
	checkScanDirectiveError(t, code, Inline, exp)
}

func TestInlineDirectiveWithUnclosedBlock(t *testing.T) {
	code := strings.Join([]string{
		"#inline c",
		"#include <stdio.h>",
	}, "\n")

	exp := strings.Join([]string{
		"example.mc:2:19: error: expect '#end-inline c' to close inline block, got EOF",
		"    2 | #include <stdio.h><EOF>",
		"      |                   ^^^^^",
	}, "\n")
	checkScanDirectiveError(t, code, Inline, exp)
}

func TestInlineDirectiveWithlockTypeUnmatched(t *testing.T) {
	code := strings.Join([]string{
		"#inline c",
		"#include <stdio.h>",
		"#end-inline asm",
		"other-code",
	}, "\n")

	exp := strings.Join([]string{
		"example.mc:4:11: error: expect '#end-inline c' to close inline block, got EOF",
		"    4 | other-code<EOF>",
		"      |           ^^^^^",
		"example.mc:3:13: note: previous possible close here",
		"    3 | #end-inline asm",
		"      |             ^^^",
		"      |             c",
	}, "\n")
	checkScanDirectiveError(t, code, Inline, exp)
}
