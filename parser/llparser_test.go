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
	expected := ast.ASTBuildDocument(
		ast.ASTBuildIncludeAngle("stdio.h"),
	)

	if err := document.EqualTo(nil, expected); err != nil {
		t.Fatalf("expected document not equal to actual:\n%s", err)
	}

	if len(document.Declarations) != 1 {
		t.Fatalf("expect 1 declaration, got %d", len(document.Declarations))
	}

	include, ok := document.Declarations[0].(*ast.PreprocessorInclude)
	if !ok {
		t.Fatalf("expect PreprocessorInclude, got %T", document.Declarations[0])
	}

	if value := include.ContentCtx.Content(); value != "stdio.h" {
		t.Fatalf("expect include content 'stdio.h', got '%s'", value)
	}
}

func TestLLParserSimplestProgram(t *testing.T) {
	code := strings.Join([]string{
		"fun main() (int) {",
		"    return 0",
		"}",
	}, "\n")

	document := runBasicTestOnCode(t, code)
	expected := ast.ASTBuildDocument(
		ast.ASTBuildFunction(
			"main",
			nil,
			ast.ASTBuildTypeList(
				ast.ASTBuildTypeListItemWithoutComma(ast.ASTBuildSimpleType("int")),
			),
			[]ast.Statement{
				ast.ASTBuildReturnStatement(
					ast.ASTBuildExpressionList(
						ast.ASTBuildExpressionListItemWithoutComma(
							ast.ASTBuildValue(0),
						),
					),
				),
			},
		),
	)
	if err := document.EqualTo(nil, expected); err != nil {
		t.Fatalf("expected document not equal to actual:\n%s", err)
	}

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

	decl, ok := document.Declarations[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("expect FunctionDeclaration, got %T", document.Declarations[0])
	}

	if len(decl.Statements) != 1 {
		t.Fatalf("expect 1 statement in function body, got %d", len(decl.Statements))
	}
}

func TestLLParserFunctionWithArguments(t *testing.T) {
	code := strings.Join([]string{
		"fun add(a int, b int) (int, int) {",
		"    return 0, 0",
		"}",
	}, "\n")

	document := runBasicTestOnCode(t, code)
	expected := ast.ASTBuildDocument(
		ast.ASTBuildFunction(
			"add",
			ast.ASTBuildArgumentList(
				ast.ASTBuildArgumentWithComma("a", ast.ASTBuildSimpleType("int")),
				ast.ASTBuildArgumentWithoutComma("b", ast.ASTBuildSimpleType("int")),
			),
			ast.ASTBuildTypeList(
				ast.ASTBuildTypeListItemWithComma(ast.ASTBuildSimpleType("int")),
				ast.ASTBuildTypeListItemWithoutComma(ast.ASTBuildSimpleType("int")),
			),
			[]ast.Statement{
				ast.ASTBuildReturnStatement(
					ast.ASTBuildExpressionList(
						ast.ASTBuildExpressionListItemWithComma(
							ast.ASTBuildValue(0),
						),
						ast.ASTBuildExpressionListItemWithoutComma(
							ast.ASTBuildValue(0),
						),
					),
				),
			},
		),
	)

	if err := document.EqualTo(nil, expected); err != nil {
		t.Fatalf("expected document not equal to actual:\n%s", err)
	}

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

	decl, ok := document.Declarations[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("expect FunctionDeclaration, got %T", document.Declarations[0])
	}

	if len(decl.Statements) != 1 {
		t.Fatalf("expect 1 statement in function body, got %d", len(decl.Statements))
	}
}

func TestLLParserReturnWithExpressionList(t *testing.T) {
	code := strings.Join([]string{
		"fun foo() (int, int, int) {",
		"    return 1, 2, 3",
		"}",
	}, "\n")

	document := runBasicTestOnCode(t, code)

	if len(document.Declarations) != 1 {
		t.Fatalf("expect 1 declaration, got %d", len(document.Declarations))
	}

	decl, ok := document.Declarations[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("expect FunctionDeclaration, got %T", document.Declarations[0])
	}

	if len(decl.Statements) != 1 {
		t.Fatalf("expect 1 statement in function body, got %d", len(decl.Statements))
	}
}

func TestLLParserReturnWithExpressionArithmetic1(t *testing.T) {
	code := strings.Join([]string{
		"fun add() (int) {",
		"    return 1 + 2",
		"}",
	}, "\n")

	document := runBasicTestOnCode(t, code)

	if len(document.Declarations) != 1 {
		t.Fatalf("expect 1 declaration, got %d", len(document.Declarations))
	}

	decl, ok := document.Declarations[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("expect FunctionDeclaration, got %T", document.Declarations[0])
	}

	if len(decl.Statements) != 1 {
		t.Fatalf("expect 1 statement in function body, got %d", len(decl.Statements))
	}
}

func TestLLParserReturnWithExpressionArithmetic2(t *testing.T) {
	code := strings.Join([]string{
		"fun add(a int, b int) (int) {",
		"    return a + b",
		"}",
	}, "\n")

	document := runBasicTestOnCode(t, code)

	if len(document.Declarations) != 1 {
		t.Fatalf("expect 1 declaration, got %d", len(document.Declarations))
	}

	decl, ok := document.Declarations[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("expect FunctionDeclaration, got %T", document.Declarations[0])
	}

	if len(decl.Statements) != 1 {
		t.Fatalf("expect 1 statement in function body, got %d", len(decl.Statements))
	}
}

func TestLLParserReturnWithExpressionArithmetic3(t *testing.T) {
	code := strings.Join([]string{
		"fun add(a int, b int) (int) {",
		"    return a + b + 3",
		"}",
	}, "\n")

	document := runBasicTestOnCode(t, code)

	if len(document.Declarations) != 1 {
		t.Fatalf("expect 1 declaration, got %d", len(document.Declarations))
	}

	decl, ok := document.Declarations[0].(*ast.FunctionDeclaration)
	if !ok {
		t.Fatalf("expect FunctionDeclaration, got %T", document.Declarations[0])
	}

	if len(decl.Statements) != 1 {
		t.Fatalf("expect 1 statement in function body, got %d", len(decl.Statements))
	}
}
