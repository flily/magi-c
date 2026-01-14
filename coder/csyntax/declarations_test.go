package csyntax

import (
	"testing"

	"strings"
)

func TestVariableDeclarationOneVariableStyle1(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("a", 0, NewIntegerLiteral(3))

	checkInterfaceCodeElement(stat)
	checkInterfaceStatement(stat)
	checkInterfaceDeclaration(stat)

	expected := "int a = 3;\n"
	checkOutputOnStyle(t, testStyle1, expected, stat)
}

func TestVariableDeclarationOneVariableStyle2(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("a", 0, NewIntegerLiteral(3))

	expected := "int a = 3;\n"
	checkOutputOnStyle(t, testStyle2, expected, stat)
}

func TestVariableDeclarationTwoVariablesStyle1(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("a", 0, NewIntegerLiteral(3))
	stat.Add("b", 0, NewIntegerLiteral(5))

	expected := "int a = 3, b = 5;\n"
	checkOutputOnStyle(t, testStyle1, expected, stat)
}

func TestVariableDeclarationTwoVariablesStyle2(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("a", 0, NewIntegerLiteral(3))
	stat.Add("b", 0, NewIntegerLiteral(5))

	expected := "int a = 3,b = 5;\n"
	checkOutputOnStyle(t, testStyle2, expected, stat)
}

func TestVariableDeclarationOnePointerVariableStyle1(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("p", 1, NewIntegerLiteral(3))

	expected := "int* p = 3;\n"
	checkOutputOnStyle(t, testStyle1, expected, stat)
}

func TestVariableDeclarationOnePointerVariableStyle2(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("p", 1, NewIntegerLiteral(3))

	expected := "int *p = 3;\n"
	checkOutputOnStyle(t, testStyle2, expected, stat)
}

func TestVariableDeclarationTwoPointerVariableStyle1(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("p", 1, NewIntegerLiteral(3))
	stat.Add("q", 2, NewIntegerLiteral(5))

	expected := "int* p = 3, ** q = 5;\n"
	checkOutputOnStyle(t, testStyle1, expected, stat)
}

func TestParameterListWrite(t *testing.T) {
	param1 := NewParameterListItem(NewType("int", 0), "a")
	param2 := NewParameterListItem(NewType("float", 1), "b")

	paramList := NewParameterList(param1, param2)

	checkInterfaceCodeElement(param1)
	checkInterfaceCodeElement(param2)
	checkInterfaceCodeElement(paramList)

	expected := "int a, float* b"
	checkOutputOnStyle(t, testStyle1, expected, paramList)
}

func TestFunctionDeclarationWithEmptyBody(t *testing.T) {
	f := NewFunctionDeclaration("add",
		NewType("int", 0),
		NewParameterList(
			NewParameterListItem(NewType("int", 0), "a"),
			NewParameterListItem(NewType("int", 0), "b"),
		),
		nil,
	)

	checkInterfaceCodeElement(f)
	checkInterfaceDeclaration(f)

	expected := strings.Join([]string{
		"int add(int a, int b) {",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle1, expected, f)
}

func TestFunctionDeclarationWithSimpleReturnStyle1(t *testing.T) {
	f := NewFunctionDeclaration("add",
		NewType("int", 0),
		NewParameterList(
			NewParameterListItem(NewType("int", 0), "a"),
			NewParameterListItem(NewType("int", 1), "b"),
		),
		nil,
	)

	returnStat := NewReturnStatement(NewIntegerLiteral(42))
	f.AddStatement(returnStat)

	expected := strings.Join([]string{
		"int add(int a, int* b) {",
		"    return 42;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle1, expected, f)
}

func TestFunctionDeclarationWithSimpleReturnStyle2(t *testing.T) {
	f := NewFunctionDeclaration("add",
		NewType("int", 0),
		NewParameterList(
			NewParameterListItem(NewType("int", 0), "a"),
			NewParameterListItem(NewType("int", 1), "b"),
		),
		nil,
	)

	returnStat := NewReturnStatement(NewIntegerLiteral(42))
	f.AddStatement(returnStat)

	expected := strings.Join([]string{
		"int add(int a,int *b)",
		"{",
		"    return 42;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle2, expected, f)
}
