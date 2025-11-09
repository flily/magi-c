package preprocessor

import (
	"testing"

	"strings"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
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
		"   1:   #include <stdio.h>",
		"        ^",
		"        here",
	}, "\n")
	checkElementContext(t, result.Hash, expHash)

	expCmd := strings.Join([]string{
		"   1:   #include <stdio.h>",
		"         ^^^^^^^",
		"         here",
	}, "\n")
	checkElementContext(t, result.Command, expCmd)

	expLBracket := strings.Join([]string{
		"   1:   #include <stdio.h>",
		"                 ^",
		"                 here",
	}, "\n")
	checkElementContext(t, result.LBracket, expLBracket)

	expContent := strings.Join([]string{
		"   1:   #include <stdio.h>",
		"                  ^^^^^^^",
		"                  here",
	}, "\n")
	checkElementContext(t, result.Content, expContent)

	expRBracket := strings.Join([]string{
		"   1:   #include <stdio.h>",
		"                         ^",
		"                         here",
	}, "\n")
	checkElementContext(t, result.RBracket, expRBracket)

	finalExp := strings.Join([]string{
		"   1:   #include <stdio.h><EOF>",
		"                          ^^^^^",
		"                          here",
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
		`   1:   #include "stdio.h"`,
		"        ^",
		"        here",
	}, "\n")
	checkElementContext(t, result.Hash, expHash)

	expCmd := strings.Join([]string{
		`   1:   #include "stdio.h"`,
		"         ^^^^^^^",
		"         here",
	}, "\n")
	checkElementContext(t, result.Command, expCmd)

	expLBracket := strings.Join([]string{
		`   1:   #include "stdio.h"`,
		"                 ^",
		"                 here",
	}, "\n")
	checkElementContext(t, result.LBracket, expLBracket)

	expContent := strings.Join([]string{
		`   1:   #include "stdio.h"`,
		"                  ^^^^^^^",
		"                  here",
	}, "\n")
	checkElementContext(t, result.Content, expContent)

	expRBracket := strings.Join([]string{
		`   1:   #include "stdio.h"`,
		"                         ^",
		"                         here",
	}, "\n")
	checkElementContext(t, result.RBracket, expRBracket)

	finalExp := strings.Join([]string{
		`   1:   #include "stdio.h"<EOF>`,
		"                          ^^^^^",
		"                          here",
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
		`   1:   #include<stdio.h>`,
		"        ^",
		"        here",
	}, "\n")
	checkElementContext(t, result.Hash, expHash)

	expCmd := strings.Join([]string{
		`   1:   #include<stdio.h>`,
		"         ^^^^^^^",
		"         here",
	}, "\n")
	checkElementContext(t, result.Command, expCmd)

	expLBracket := strings.Join([]string{
		`   1:   #include<stdio.h>`,
		"                ^",
		"                here",
	}, "\n")
	checkElementContext(t, result.LBracket, expLBracket)

	expContent := strings.Join([]string{
		`   1:   #include<stdio.h>`,
		"                 ^^^^^^^",
		"                 here",
	}, "\n")
	checkElementContext(t, result.Content, expContent)

	expRBracket := strings.Join([]string{
		`   1:   #include<stdio.h>`,
		"                        ^",
		"                        here",
	}, "\n")
	checkElementContext(t, result.RBracket, expRBracket)

	finalExp := strings.Join([]string{
		`   1:   #include<stdio.h><EOF>`,
		"                         ^^^^^",
		"                         here",
	}, "\n")
	checkElementContext(t, final, finalExp)
}

func TestIncludeDirectiveWithWrongQuote(t *testing.T) {
	code := strings.Join([]string{
		`#include stdio.h`,
	}, "\n")

	exp := strings.Join([]string{
		`   1:   #include stdio.h`,
		"                 ^",
		"                 expected '<' or '\"' after '#include', got 's'",
	}, "\n")
	checkScanDirectiveError(t, code, Include, exp)
}

func TestIncludeDirectiveWithNoName(t *testing.T) {
	code := strings.Join([]string{
		`#include <>`,
	}, "\n")

	exp := strings.Join([]string{
		`   1:   #include <>`,
		"                 ^",
		"                 expected file name after '#include', got empty string",
	}, "\n")
	checkScanDirectiveError(t, code, Include, exp)
}

func TestIncludeDirectiveWithUnclosedQuote(t *testing.T) {
	code := strings.Join([]string{
		`#include <stdio.h`,
		"",
	}, "\n")

	cursor := context.NewCursorFromString("example.c", code)
	_, err := scanDirectiveOn(cursor, Include)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	exp := strings.Join([]string{
		`   1:   #include <stdio.h<EOL LF>`,
		"                 ^       ^^^^^^^^",
		"                 quote not closed",
	}, "\n")
	got := err.Error()
	if got != exp {
		t.Errorf("expected error message:\n%s\ngot:\n%s", exp, got)
	}
}

func TestIncludeDirectiveWithQuoteNotMatched(t *testing.T) {
	code := strings.Join([]string{
		`#include "stdio.h>`,
	}, "\n")

	cursor := context.NewCursorFromString("example.c", code)
	_, err := scanDirectiveOn(cursor, Include)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	exp := strings.Join([]string{
		`   1:   #include "stdio.h>`,
		"                 ^       ^",
		"                 quote mismatch, expected '\"', got '>'",
	}, "\n")
	got := err.Error()
	if got != exp {
		t.Errorf("expected error message:\n%s\ngot:\n%s", exp, got)
	}
}
