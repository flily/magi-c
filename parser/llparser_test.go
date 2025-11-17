package parser

import (
	"strings"
	"testing"

	"github.com/flily/magi-c/ast"
)

func TestLLParserSimpleStatement(t *testing.T) {
	code := strings.Join([]string{
		"#include <stdio.h>",
	}, "\n")

	document := runBasicTestOnCode(t, code)

	if len(document.Declarations) != 1 {
		t.Fatalf("expect 1 declaration, got %d", len(document.Declarations))
	}

	include, ok := document.Declarations[0].(*ast.PreprocessorInclude)
	if !ok {
		t.Fatalf("expect PreprocessorInclude, got %T", document.Declarations[0])
	}

	if value := include.Content.Content(); value != "stdio.h" {
		t.Fatalf("expect include content 'stdio.h', got '%s'", value)
	}
}

func TestLLParserSimplestProgram(t *testing.T) {
	code := strings.Join([]string{
		"fun main() () {",
		"    return 0",
		"}",
	}, "\n")

	document := runBasicTestOnCode(t, code)

	if len(document.Declarations) != 1 {
		t.Fatalf("expect 1 declaration, got %d", len(document.Declarations))
	}

	main, ok := document.Declarations[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("expect FunctionDeclaration, got %T", document.Declarations[0])
	}

	if main.Name.Name != "main" {
		t.Fatalf("expect function name 'main', got '%s'", main.Name.Name)
	}
}

func TestLLParserFunctionWithArguments(t *testing.T) {
	code := strings.Join([]string{
		"fun add(a int, b int) () {",
		"    return 0",
		"}",
	}, "\n")

	document := runBasicTestOnCode(t, code)

	if len(document.Declarations) != 1 {
		t.Fatalf("expect 1 declaration, got %d", len(document.Declarations))
	}

	add, ok := document.Declarations[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("expect FunctionDeclaration, got %T", document.Declarations[0])
	}

	if add.Name.Name != "add" {
		t.Fatalf("expect function name 'add', got '%s'", add.Name.Name)
	}
}
