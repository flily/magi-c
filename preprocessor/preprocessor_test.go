package preprocessor

import (
	"testing"

	"strings"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
)

func testScanDirectiveCorrect(t *testing.T, code string, prep PreprocessorInitializer) (ast.Node, *context.Context) {
	cursor := context.NewCursorFromString("example.c", code)
	node, err := scanDirectiveOn(cursor, prep)
	if err != nil {
		t.Fatalf("unexpected error when scan directive:\n%v", err)
	}

	if node == nil {
		t.Fatalf("expect non-nil node when scan directive")
	}

	_, final := cursor.CurrentChar()
	return node, final
}

func checkScanDirectiveError(t *testing.T, code string, prep PreprocessorInitializer, expected string) {
	cursor := context.NewCursorFromString("example.c", code)
	_, err := scanDirectiveOn(cursor, prep)
	if err == nil {
		t.Fatalf("expect error when scan directive")
	}

	got := err.Error()
	if got != expected {
		t.Errorf("expect error message:\n%s\ngot:\n%s", expected, got)
	}
}

func checkElementContext(t *testing.T, ctx *context.Context, expected string) {
	got := ctx.HighlightText("here")
	if got != expected {
		t.Errorf("expect highlight on element:\n%s\ngot:\n%s", expected, got)
	}
}

func TestScanDirective(t *testing.T) {
	code := strings.Join([]string{
		"#include <stdio.h>",
	}, "\n")

	cursor := context.NewCursorFromString("example.c", code)
	cmd, hashCtx, cmdCtx, err := ScanDirective(cursor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cmd != "include" {
		t.Errorf("expect command 'include', got '%s'", cmd)
	}

	hashExp := strings.Join([]string{
		"   1:   #include <stdio.h>",
		"        ^",
		"        here",
	}, "\n")
	checkElementContext(t, hashCtx, hashExp)

	cmdExp := strings.Join([]string{
		"   1:   #include <stdio.h>",
		"         ^^^^^^^",
		"         here",
	}, "\n")
	checkElementContext(t, cmdCtx, cmdExp)
}

func TestScanDirectiveStartWithNoHash(t *testing.T) {
	code := strings.Join([]string{
		"include <stdio.h>",
	}, "\n")

	exp := strings.Join([]string{
		"   1:   include <stdio.h>",
		"        ^",
		"        expect '#' at the beginning of preprocessor directive, got 'i'",
	}, "\n")
	checkScanDirectiveError(t, code, Include, exp)
}

func TestScanDirectiveWithHashIsNotTheFirstChar(t *testing.T) {
	code := strings.Join([]string{
		"while #include <stdio.h>",
	}, "\n")

	cursor := context.NewCursorFromString("example.c", code)
	cursor.Skip(6)

	hashChar, current := cursor.CurrentChar()
	if hashChar != '#' {
		t.Fatalf("expect current char to be '#', got '%c'", hashChar)
	}

	expPosition := strings.Join([]string{
		"   1:   while #include <stdio.h>",
		"              ^",
		"              here",
	}, "\n")
	checkElementContext(t, current, expPosition)

	expError := strings.Join([]string{
		"   1:   while #include <stdio.h>",
		"              ^",
		"              '#' must be the first non-whitespace character in the line",
	}, "\n")
	_, err := scanDirectiveOn(cursor, Include)
	if err == nil {
		t.Fatalf("expect error when scan directive")
	}
	gotError := err.Error()
	if gotError != expError {
		t.Errorf("expect error message:\n%s\ngot:\n%s", expError, gotError)
	}
}

func TestScanDirectiveWithNoName(t *testing.T) {
	code := strings.Join([]string{
		"#  <stdio.h>",
	}, "\n")

	exp := strings.Join([]string{
		"   1:   #  <stdio.h>",
		"         ^",
		"         expect preprocessor directive name after '#', got empty string",
	}, "\n")
	checkScanDirectiveError(t, code, Include, exp)
}
