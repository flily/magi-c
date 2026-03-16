package csyntax

import (
	"testing"

	"strings"
)

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
	checkInterfaceDefinition(f)

	expected := strings.Join([]string{
		"int add(int a, int b)",
		"{",
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
		"int add(int a, int* b)",
		"{",
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
