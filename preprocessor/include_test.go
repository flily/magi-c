package preprocessor

import (
	"testing"

	"strings"

	"github.com/flily/magi-c/ast"
)

func TestIncludeDirectiveAngleQuote(t *testing.T) {
	code := strings.Join([]string{
		"#include <stdio.h>",
	}, "\n")

	node, final := testScanDirectiveCorrect(t, code, Include)
	result, ok := node.(*ast.PreprocessorInclude)
	if !ok {
		t.Fatalf("expected PreprocessorInclude node, got %T", node)
	}

	expHash := strings.Join([]string{
		"    1 | #include <stdio.h>",
		"      | ^",
		"      | here",
	}, "\n")
	checkElementContext(t, result.Hash, expHash)

	expCmd := strings.Join([]string{
		"    1 | #include <stdio.h>",
		"      |  ^^^^^^^",
		"      |  here",
	}, "\n")
	checkElementContext(t, result.Command, expCmd)

	expLBracket := strings.Join([]string{
		"    1 | #include <stdio.h>",
		"      |          ^",
		"      |          here",
	}, "\n")
	checkElementContext(t, result.LBracketCtx, expLBracket)

	expContent := strings.Join([]string{
		"    1 | #include <stdio.h>",
		"      |           ^^^^^^^",
		"      |           here",
	}, "\n")
	checkElementContext(t, result.ContentCtx, expContent)

	expRBracket := strings.Join([]string{
		"    1 | #include <stdio.h>",
		"      |                  ^",
		"      |                  here",
	}, "\n")
	checkElementContext(t, result.RBracketCtx, expRBracket)

	finalExp := strings.Join([]string{
		"    1 | #include <stdio.h><EOF>",
		"      |                   ^^^^^",
		"      |                   here",
	}, "\n")
	checkElementContext(t, final, finalExp)
}

func TestIncludeDirectiveDoubleQuote(t *testing.T) {
	code := strings.Join([]string{
		`#include "stdio.h"`,
	}, "\n")

	node, final := testScanDirectiveCorrect(t, code, Include)
	result, ok := node.(*ast.PreprocessorInclude)
	if !ok {
		t.Fatalf("expected PreprocessorInclude node, got %T", node)
	}

	expHash := strings.Join([]string{
		`    1 | #include "stdio.h"`,
		"      | ^",
		"      | here",
	}, "\n")
	checkElementContext(t, result.Hash, expHash)

	expCmd := strings.Join([]string{
		`    1 | #include "stdio.h"`,
		"      |  ^^^^^^^",
		"      |  here",
	}, "\n")
	checkElementContext(t, result.Command, expCmd)

	expLBracket := strings.Join([]string{
		`    1 | #include "stdio.h"`,
		"      |          ^",
		"      |          here",
	}, "\n")
	checkElementContext(t, result.LBracketCtx, expLBracket)

	expContent := strings.Join([]string{
		`    1 | #include "stdio.h"`,
		"      |           ^^^^^^^",
		"      |           here",
	}, "\n")
	checkElementContext(t, result.ContentCtx, expContent)

	expRBracket := strings.Join([]string{
		`    1 | #include "stdio.h"`,
		"      |                  ^",
		"      |                  here",
	}, "\n")
	checkElementContext(t, result.RBracketCtx, expRBracket)

	finalExp := strings.Join([]string{
		`    1 | #include "stdio.h"<EOF>`,
		"      |                   ^^^^^",
		"      |                   here",
	}, "\n")
	checkElementContext(t, final, finalExp)
}

func TestIncludeDirectiveWithoutSpace(t *testing.T) {
	code := strings.Join([]string{
		`#include<stdio.h>`,
	}, "\n")

	node, final := testScanDirectiveCorrect(t, code, Include)
	result, ok := node.(*ast.PreprocessorInclude)
	if !ok {
		t.Fatalf("expected PreprocessorInclude node, got %T", node)
	}

	expHash := strings.Join([]string{
		`    1 | #include<stdio.h>`,
		"      | ^",
		"      | here",
	}, "\n")
	checkElementContext(t, result.Hash, expHash)

	expCmd := strings.Join([]string{
		`    1 | #include<stdio.h>`,
		"      |  ^^^^^^^",
		"      |  here",
	}, "\n")
	checkElementContext(t, result.Command, expCmd)

	expLBracket := strings.Join([]string{
		`    1 | #include<stdio.h>`,
		"      |         ^",
		"      |         here",
	}, "\n")
	checkElementContext(t, result.LBracketCtx, expLBracket)

	expContent := strings.Join([]string{
		`    1 | #include<stdio.h>`,
		"      |          ^^^^^^^",
		"      |          here",
	}, "\n")
	checkElementContext(t, result.ContentCtx, expContent)

	expRBracket := strings.Join([]string{
		`    1 | #include<stdio.h>`,
		"      |                 ^",
		"      |                 here",
	}, "\n")
	checkElementContext(t, result.RBracketCtx, expRBracket)

	finalExp := strings.Join([]string{
		`    1 | #include<stdio.h><EOF>`,
		"      |                  ^^^^^",
		"      |                  here",
	}, "\n")
	checkElementContext(t, final, finalExp)
}

func TestIncludeDirectiveWithWrongQuote(t *testing.T) {
	code := strings.Join([]string{
		`#include stdio.h`,
	}, "\n")

	exp := strings.Join([]string{
		"example.mc:1:10: error: expected '<' or '\"' after '#include', got 's'",
		`    1 | #include stdio.h`,
		"      |          ^",
		"      |          < or \"",
	}, "\n")
	checkScanDirectiveError(t, code, Include, exp)
}

func TestIncludeDirectiveWithNoName(t *testing.T) {
	code := strings.Join([]string{
		`#include <>`,
	}, "\n")

	exp := strings.Join([]string{
		"example.mc:1:11: error: expected file name after '#include', got empty string",
		`    1 | #include <>`,
		"      |           ^",
	}, "\n")
	checkScanDirectiveError(t, code, Include, exp)
}

func TestIncludeDirectiveWithUnclosedQuote(t *testing.T) {
	code := strings.Join([]string{
		`#include <stdio.h`,
		"",
	}, "\n")

	exp := strings.Join([]string{
		"example.mc:1:18: error: quote not closed",
		`    1 | #include <stdio.h<EOL LF>`,
		"      |                  ^^^^^^^^",
		"      |                  >",
	}, "\n")
	checkScanDirectiveError(t, code, Include, exp)
}

func TestIncludeDirectiveWithQuoteNotMatched(t *testing.T) {
	code := strings.Join([]string{
		`#include "stdio.h>`,
	}, "\n")

	exp := strings.Join([]string{
		"example.mc:1:18: error: quote mismatch, expected '\"', got '>'",
		`    1 | #include "stdio.h>`,
		"      |                  ^",
		"      |                  \"",
		"example.mc:1:10: note: previous quote is '\"'",
		`    1 | #include "stdio.h>`,
		"      |          ^",
	}, "\n")
	checkScanDirectiveError(t, code, Include, exp)
}
