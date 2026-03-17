package csyntax

import (
	"testing"
)

func TestVariableDeclarationOneVariableStyle1(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("a", 0, NewIntegerLiteral(3))

	checkInterfaceCodeElement(decl)
	checkInterfaceDeclaration(decl)

	expected := "int a = 3"
	checkOutputOnStyle(t, testStyle1, expected, decl)
}

func TestVariableDeclarationOneVariableStyle2(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("a", 0, NewIntegerLiteral(3))

	checkInterfaceCodeElement(decl)
	checkInterfaceDeclaration(decl)

	expected := "int a = 3"
	checkOutputOnStyle(t, testStyle2, expected, decl)
}

func TestVariableDeclarationTwoVariablesStyle1(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("a", 0, NewIntegerLiteral(3))
	decl.Add("b", 0, NewIntegerLiteral(5))

	checkInterfaceCodeElement(decl)
	checkInterfaceDeclaration(decl)

	expected := "int a = 3, b = 5"
	checkOutputOnStyle(t, testStyle1, expected, decl)
}

func TestVariableDeclarationTwoVariablesStyle2(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("a", 0, NewIntegerLiteral(3))
	decl.Add("b", 0, NewIntegerLiteral(5))

	checkInterfaceCodeElement(decl)
	checkInterfaceDeclaration(decl)

	expected := "int a = 3,b = 5"
	checkOutputOnStyle(t, testStyle2, expected, decl)
}

func TestVariableDeclarationOnePointerVariableStyle1(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("p", 1, NewIntegerLiteral(3))

	checkInterfaceCodeElement(decl)
	checkInterfaceDeclaration(decl)

	expected := "int* p = 3"
	checkOutputOnStyle(t, testStyle1, expected, decl)
}

func TestVariableDeclarationOnePointerVariableStyle2(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("p", 1, NewIntegerLiteral(3))

	checkInterfaceCodeElement(decl)
	checkInterfaceDeclaration(decl)

	expected := "int *p = 3"
	checkOutputOnStyle(t, testStyle2, expected, decl)
}

func TestVariableDeclarationTwoPointerVariableStyle1(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("p", 1, NewIntegerLiteral(3))
	decl.Add("q", 2, NewIntegerLiteral(5))

	checkInterfaceCodeElement(decl)
	checkInterfaceDeclaration(decl)

	expected := "int* p = 3, ** q = 5"
	checkOutputOnStyle(t, testStyle1, expected, decl)
}

func TestVariableDeclarationTwoPointerVariableStyle2(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("p", 1, NewIntegerLiteral(3))
	decl.Add("q", 2, NewIntegerLiteral(5))

	checkInterfaceCodeElement(decl)
	checkInterfaceDeclaration(decl)

	expected := "int *p = 3, **q = 5"
	checkOutputOnStyle(t, testStyle2, expected, decl)
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
