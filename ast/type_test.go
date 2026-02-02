package ast

import (
	"testing"

	"strings"
)

func TestSimpleType(t *testing.T) {
	text := "* * lorem"
	ctxList := generateTestWords(text)

	asterisks := []*TerminalToken{
		NewTerminalToken(ctxList[0], Asterisk),
		NewTerminalToken(ctxList[1], Asterisk),
	}
	identifier := NewIdentifier(ctxList[2])
	simpleType := NewSimpleType(asterisks, identifier)

	var _ Type = simpleType
	simpleType.typeNode()

	expected := ASTBuildSimpleType("**lorem")

	if err := simpleType.EqualTo(nil, expected); err != nil {
		t.Fatalf("SimpleType not equal: %s", err)
	}
}

func TestSimpleTypeNotEqualOnType(t *testing.T) {
	text := "* * lorem"
	ctxList := generateTestWords(text)

	asterisks := []*TerminalToken{
		NewTerminalToken(ctxList[0], Asterisk),
		NewTerminalToken(ctxList[1], Asterisk),
	}
	identifier := NewIdentifier(ctxList[2])
	simpleType := NewSimpleType(asterisks, identifier)

	expected := ASTBuildValue(42)
	message := strings.Join([]string{
		"test.txt:1:1: error: expect a *ast.IntegerLiteral, got a *ast.SimpleType",
		"    1 | * * lorem",
		"      | ^ ^ ^^^^^",
		"      | *ast.IntegerLiteral",
	}, "\n")

	err := simpleType.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestSimpleTypeNotEqualOnAsterisks(t *testing.T) {
	text := "* * lorem"
	ctxList := generateTestWords(text)

	asterisks := []*TerminalToken{
		NewTerminalToken(ctxList[0], Asterisk),
		NewTerminalToken(ctxList[1], Asterisk),
	}
	identifier := NewIdentifier(ctxList[2])
	simpleType := NewSimpleType(asterisks, identifier)

	expected := ASTBuildSimpleType("*lorem")
	message := strings.Join([]string{
		"test.txt:1:1: error: wrong number of POINTER ASTERISK: expected 1, got 2",
		"    1 | * * lorem",
		"      | ^ ^",
	}, "\n")

	err := simpleType.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestSimpleTypeNotEqualOnIdentifier(t *testing.T) {
	text := "* * lorem"
	ctxList := generateTestWords(text)

	asterisks := []*TerminalToken{
		NewTerminalToken(ctxList[0], Asterisk),
		NewTerminalToken(ctxList[1], Asterisk),
	}
	identifier := NewIdentifier(ctxList[2])
	simpleType := NewSimpleType(asterisks, identifier)

	expected := ASTBuildSimpleType("**ipsum")
	message := strings.Join([]string{
		"test.txt:1:5: error: wrong identifier name, expect 'ipsum', got 'lorem'",
		"    1 | * * lorem",
		"      |     ^^^^^",
		"      |     ipsum",
	}, "\n")

	err := simpleType.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestArgumentDeclaration(t *testing.T) {
	text := "argName * int ,"
	ctxList := generateTestWords(text)

	argName := NewIdentifier(ctxList[0])
	argType := NewSimpleType([]*TerminalToken{NewTerminalToken(ctxList[1], Asterisk)}, NewIdentifier(ctxList[2]))
	argComma := NewTerminalToken(ctxList[3], Comma)
	argDecl := NewArgumentDeclaration(argName, argType, argComma)

	expected := ASTBuildArgumentWithComma("argName", "*int")

	if err := argDecl.EqualTo(nil, expected); err != nil {
		t.Fatalf("ArgumentDeclaration not equal: %s", err)
	}
}

func TestArgumentDeclarationNotEqualOnNodeType(t *testing.T) {
	text := "lorem * int ,"
	ctxList := generateTestWords(text)

	argName := NewIdentifier(ctxList[0])
	argType := NewSimpleType([]*TerminalToken{NewTerminalToken(ctxList[1], Asterisk)}, NewIdentifier(ctxList[2]))
	argComma := NewTerminalToken(ctxList[3], Comma)
	argDecl := NewArgumentDeclaration(argName, argType, argComma)

	expected := ASTBuildValue(42)
	message := strings.Join([]string{
		"test.txt:1:1: error: expect a *ast.IntegerLiteral, got a *ast.ArgumentDeclaration",
		"    1 | lorem * int ,",
		"      | ^^^^^ ^ ^^^ ^",
		"      | *ast.IntegerLiteral",
	}, "\n")

	err := argDecl.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestArgumentDeclarationNotEqualOnName(t *testing.T) {
	text := "lorem * int ,"
	ctxList := generateTestWords(text)

	argName := NewIdentifier(ctxList[0])
	argType := NewSimpleType([]*TerminalToken{NewTerminalToken(ctxList[1], Asterisk)}, NewIdentifier(ctxList[2]))
	argComma := NewTerminalToken(ctxList[3], Comma)
	argDecl := NewArgumentDeclaration(argName, argType, argComma)

	expected := ASTBuildArgumentWithComma("ipsum", "*int")
	message := strings.Join([]string{
		"test.txt:1:1: error: wrong identifier name, expect 'ipsum', got 'lorem'",
		"    1 | lorem * int ,",
		"      | ^^^^^",
		"      | ipsum",
	}, "\n")

	err := argDecl.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestArgumentDeclarationNotEqualOnType(t *testing.T) {
	text := "lorem * int ,"
	ctxList := generateTestWords(text)

	argName := NewIdentifier(ctxList[0])
	argType := NewSimpleType([]*TerminalToken{NewTerminalToken(ctxList[1], Asterisk)}, NewIdentifier(ctxList[2]))
	argComma := NewTerminalToken(ctxList[3], Comma)
	argDecl := NewArgumentDeclaration(argName, argType, argComma)

	expected := ASTBuildArgumentWithComma("lorem", "int")
	message := strings.Join([]string{
		"test.txt:1:7: error: wrong number of POINTER ASTERISK: expected 0, got 1",
		"    1 | lorem * int ,",
		"      |       ^",
	}, "\n")

	err := argDecl.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestArgumentDeclarationNotEqualOnComma(t *testing.T) {
	text := "lorem * int ,"
	ctxList := generateTestWords(text)

	argName := NewIdentifier(ctxList[0])
	argType := NewSimpleType([]*TerminalToken{NewTerminalToken(ctxList[1], Asterisk)}, NewIdentifier(ctxList[2]))
	argComma := NewTerminalToken(ctxList[3], Comma)
	argDecl := NewArgumentDeclaration(argName, argType, argComma)

	expected := ASTBuildArgumentWithoutComma("lorem", "*int")
	message := strings.Join([]string{
		"test.txt:1:13: error: unexpected *ast.TerminalToken found",
		"    1 | lorem * int ,",
		"      |             ^",
		"      |             unexpected token",
	}, "\n")

	err := argDecl.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestArgumentList(t *testing.T) {
	text := "a int , b float , c string"
	ctxList := generateTestWords(text)

	args := NewArgumentList()
	args.Add(NewIdentifier(ctxList[0]), NewSimpleType(nil, NewIdentifier(ctxList[1])), NewTerminalToken(ctxList[2], Comma))
	args.Add(NewIdentifier(ctxList[3]), NewSimpleType(nil, NewIdentifier(ctxList[4])), NewTerminalToken(ctxList[5], Comma))
	args.Add(NewIdentifier(ctxList[6]), NewSimpleType(nil, NewIdentifier(ctxList[7])), nil)

	expected := ASTBuildArgumentList(
		ASTBuildArgumentWithComma("a", "int"),
		ASTBuildArgumentWithComma("b", "float"),
		ASTBuildArgumentWithoutComma("c", "string"),
	)

	if err := args.EqualTo(nil, expected); err != nil {
		t.Fatalf("ArgumentList not equal: %s", err)
	}

	if args.Length() != 3 {
		t.Fatalf("ArgumentList length expect 3, got %d", args.Length())
	}
}

func TestArgumentListNotEqualOnNodeType(t *testing.T) {
	text := "a int , b float , c string"
	ctxList := generateTestWords(text)

	args := NewArgumentList()
	args.Add(NewIdentifier(ctxList[0]), NewSimpleType(nil, NewIdentifier(ctxList[1])), NewTerminalToken(ctxList[2], Comma))
	args.Add(NewIdentifier(ctxList[3]), NewSimpleType(nil, NewIdentifier(ctxList[4])), NewTerminalToken(ctxList[5], Comma))
	args.Add(NewIdentifier(ctxList[6]), NewSimpleType(nil, NewIdentifier(ctxList[7])), nil)

	expected := ASTBuildValue(42)
	message := strings.Join([]string{
		"test.txt:1:1: error: expect a *ast.IntegerLiteral, got a *ast.ArgumentList",
		"    1 | a int , b float , c string",
		"      | ^ ^^^ ^ ^ ^^^^^ ^ ^ ^^^^^^",
		"      | *ast.IntegerLiteral",
	}, "\n")

	err := args.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestTypeListItem(t *testing.T) {
	text := "* * int ,"
	ctxList := generateTestWords(text)

	typeItem := NewTypeListItem(
		NewSimpleType(
			[]*TerminalToken{
				NewTerminalToken(ctxList[0], Asterisk),
				NewTerminalToken(ctxList[1], Asterisk),
			},
			NewIdentifier(ctxList[2]),
		),
		NewTerminalToken(ctxList[3], Comma),
	)

	expected := ASTBuildTypeListItemWithComma("**int")

	if err := typeItem.EqualTo(nil, expected); err != nil {
		t.Fatalf("TypeListItem not equal: %s", err)
	}
}

func TestTypeListItemNotEqualOnNodeType(t *testing.T) {
	text := "* * int ,"
	ctxList := generateTestWords(text)

	typeItem := NewTypeListItem(
		NewSimpleType(
			[]*TerminalToken{
				NewTerminalToken(ctxList[0], Asterisk),
				NewTerminalToken(ctxList[1], Asterisk),
			},
			NewIdentifier(ctxList[2]),
		),
		NewTerminalToken(ctxList[3], Comma),
	)

	expected := ASTBuildValue(42)
	message := strings.Join([]string{
		"test.txt:1:1: error: expect a *ast.IntegerLiteral, got a *ast.TypeListItem",
		"    1 | * * int ,",
		"      | ^ ^ ^^^ ^",
		"      | *ast.IntegerLiteral",
	}, "\n")

	err := typeItem.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestTypeListItemNotEqualOnType(t *testing.T) {
	text := "* * int ,"
	ctxList := generateTestWords(text)

	typeItem := NewTypeListItem(
		NewSimpleType(
			[]*TerminalToken{
				NewTerminalToken(ctxList[0], Asterisk),
				NewTerminalToken(ctxList[1], Asterisk),
			},
			NewIdentifier(ctxList[2]),
		),
		NewTerminalToken(ctxList[3], Comma),
	)

	expected := ASTBuildTypeListItemWithComma("*int")
	message := strings.Join([]string{
		"test.txt:1:1: error: wrong number of POINTER ASTERISK: expected 1, got 2",
		"    1 | * * int ,",
		"      | ^ ^",
	}, "\n")

	err := typeItem.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestTypeListItemNotEqualOnComma(t *testing.T) {
	text := "* * int ,"
	ctxList := generateTestWords(text)

	typeItem := NewTypeListItem(
		NewSimpleType(
			[]*TerminalToken{
				NewTerminalToken(ctxList[0], Asterisk),
				NewTerminalToken(ctxList[1], Asterisk),
			},
			NewIdentifier(ctxList[2]),
		),
		NewTerminalToken(ctxList[3], Comma),
	)

	expected := ASTBuildTypeListItemWithoutComma("**int")
	message := strings.Join([]string{
		"test.txt:1:9: error: unexpected *ast.TerminalToken found",
		"    1 | * * int ,",
		"      |         ^",
		"      |         unexpected token",
	}, "\n")

	err := typeItem.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestTypeList(t *testing.T) {
	text := "* int , * * float , string"
	ctxList := generateTestWords(text)

	types := NewTypeList()
	types.Add(NewSimpleType([]*TerminalToken{NewTerminalToken(ctxList[0], Asterisk)}, NewIdentifier(ctxList[1])), NewTerminalToken(ctxList[2], Comma))
	types.Add(NewSimpleType([]*TerminalToken{NewTerminalToken(ctxList[3], Asterisk), NewTerminalToken(ctxList[4], Asterisk)}, NewIdentifier(ctxList[5])), NewTerminalToken(ctxList[6], Comma))
	types.Add(NewSimpleType(nil, NewIdentifier(ctxList[7])), nil)

	expected := ASTBuildTypeList(
		ASTBuildTypeListItemWithComma("*int"),
		ASTBuildTypeListItemWithComma("**float"),
		ASTBuildTypeListItemWithoutComma("string"),
	)

	if err := types.EqualTo(nil, expected); err != nil {
		t.Fatalf("TypeList not equal: %s", err)
	}

	if types.Length() != 3 {
		t.Fatalf("TypeList length expect 3, got %d", types.Length())
	}
}

func TestTypeListNotEqualOnNodeType(t *testing.T) {
	text := "* int , * * float , string"
	ctxList := generateTestWords(text)

	types := NewTypeList()
	types.Add(NewSimpleType([]*TerminalToken{NewTerminalToken(ctxList[0], Asterisk)}, NewIdentifier(ctxList[1])), NewTerminalToken(ctxList[2], Comma))
	types.Add(NewSimpleType([]*TerminalToken{NewTerminalToken(ctxList[3], Asterisk), NewTerminalToken(ctxList[4], Asterisk)}, NewIdentifier(ctxList[5])), NewTerminalToken(ctxList[6], Comma))
	types.Add(NewSimpleType(nil, NewIdentifier(ctxList[7])), nil)

	expected := ASTBuildValue(42)
	message := strings.Join([]string{
		"test.txt:1:1: error: expect a *ast.IntegerLiteral, got a *ast.TypeList",
		"    1 | * int , * * float , string",
		"      | ^ ^^^ ^ ^ ^ ^^^^^ ^ ^^^^^^",
		"      | *ast.IntegerLiteral",
	}, "\n")

	err := types.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestTypeListNotEqualOnSize(t *testing.T) {
	text := "* int , * * float , string"
	ctxList := generateTestWords(text)

	types := NewTypeList()
	types.Add(NewSimpleType([]*TerminalToken{NewTerminalToken(ctxList[0], Asterisk)}, NewIdentifier(ctxList[1])), NewTerminalToken(ctxList[2], Comma))
	types.Add(NewSimpleType([]*TerminalToken{NewTerminalToken(ctxList[3], Asterisk), NewTerminalToken(ctxList[4], Asterisk)}, NewIdentifier(ctxList[5])), NewTerminalToken(ctxList[6], Comma))
	types.Add(NewSimpleType(nil, NewIdentifier(ctxList[7])), nil)

	expected := ASTBuildTypeList(
		ASTBuildTypeListItemWithComma("*int"),
		ASTBuildTypeListItemWithComma("**float"),
	)
	message := strings.Join([]string{
		"test.txt:1:1: error: wrong number of TYPE LIST: expected 2, got 3",
		"    1 | * int , * * float , string",
		"      | ^ ^^^ ^ ^ ^ ^^^^^ ^ ^^^^^^",
	}, "\n")

	err := types.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}
