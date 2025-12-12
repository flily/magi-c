package ast

import (
	"testing"

	"strings"
)

func TestArgumentList(t *testing.T) {
	text := "a int , b float , c string"
	ctxList := generateTestWords(text)

	args := NewArgumentList()
	args.Add(NewIdentifier(ctxList[0]), NewSimpleType(nil, NewIdentifier(ctxList[1])), NewTerminalToken(ctxList[2], Comma))
	args.Add(NewIdentifier(ctxList[3]), NewSimpleType(nil, NewIdentifier(ctxList[4])), NewTerminalToken(ctxList[5], Comma))
	args.Add(NewIdentifier(ctxList[6]), NewSimpleType(nil, NewIdentifier(ctxList[7])), nil)

	expected := ASTBuildArgumentList(
		ASTBuildArgumentWithComma("a", ASTBuildSimpleType("int")),
		ASTBuildArgumentWithComma("b", ASTBuildSimpleType("float")),
		ASTBuildArgumentWithoutComma("c", ASTBuildSimpleType("string")),
	)

	if err := args.EqualTo(nil, expected); err != nil {
		t.Fatalf("ArgumentList not equal: %s", err)
	}
}

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
		"   1:   * * lorem",
		"        ^ ^ ^^^^^",
		"        expect a *ast.IntegerLiteral",
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
		"   1:   * * lorem",
		"        ^ ^",
		"        wrong number of POINTER ASTERISK: expected 1, got 2",
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
		"   1:   * * lorem",
		"            ^^^^^",
		"            wrong identifier name, expect 'ipsum', got 'lorem'",
	}, "\n")

	err := simpleType.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("expect a error but got nil")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}
