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

	builder, writer := makeTestWriter(testStyle1)
	err := stat.Write(writer, 0)
	if err != nil {
		t.Fatalf("VariableDeclaration Write failed: %s", err)
	}

	expected := "int a = 3;\n"
	result := builder.String()
	if result != expected {
		t.Fatalf("VariableDeclaration Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestVariableDeclarationOneVariableStyle2(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("a", 0, NewIntegerLiteral(3))

	builder, writer := makeTestWriter(testStyle2)
	err := stat.Write(writer, 0)
	if err != nil {
		t.Fatalf("VariableDeclaration Write failed: %s", err)
	}

	expected := "int a = 3;\n"
	result := builder.String()
	if result != expected {
		t.Fatalf("VariableDeclaration Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestVariableDeclarationTwoVariablesStyle1(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("a", 0, NewIntegerLiteral(3))
	stat.Add("b", 0, NewIntegerLiteral(5))

	builder, writer := makeTestWriter(testStyle1)
	err := stat.Write(writer, 0)
	if err != nil {
		t.Fatalf("VariableDeclaration Write failed: %s", err)
	}

	expected := "int a = 3, b = 5;\n"
	result := builder.String()
	if result != expected {
		t.Fatalf("VariableDeclaration Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestVariableDeclarationTwoVariablesStyle2(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("a", 0, NewIntegerLiteral(3))
	stat.Add("b", 0, NewIntegerLiteral(5))

	builder, writer := makeTestWriter(testStyle2)
	err := stat.Write(writer, 0)
	if err != nil {
		t.Fatalf("VariableDeclaration Write failed: %s", err)
	}

	expected := "int a = 3,b = 5;\n"
	result := builder.String()
	if result != expected {
		t.Fatalf("VariableDeclaration Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestVariableDeclarationOnePointerVariableStyle1(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("p", 1, NewIntegerLiteral(3))

	builder, writer := makeTestWriter(testStyle1)
	err := stat.Write(writer, 0)
	if err != nil {
		t.Fatalf("VariableDeclaration Write failed: %s", err)
	}

	expected := "int* p = 3;\n"
	result := builder.String()
	if result != expected {
		t.Fatalf("VariableDeclaration Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestVariableDeclarationOnePointerVariableStyle2(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("p", 1, NewIntegerLiteral(3))

	builder, writer := makeTestWriter(testStyle2)
	err := stat.Write(writer, 0)
	if err != nil {
		t.Fatalf("VariableDeclaration Write failed: %s", err)
	}

	expected := "int *p = 3;\n"
	result := builder.String()
	if result != expected {
		t.Fatalf("VariableDeclaration Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestVariableDeclarationTwoPointerVariableStyle1(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("p", 1, NewIntegerLiteral(3))
	stat.Add("q", 2, NewIntegerLiteral(5))

	builder, writer := makeTestWriter(testStyle1)
	err := stat.Write(writer, 0)
	if err != nil {
		t.Fatalf("VariableDeclaration Write failed: %s", err)
	}

	expected := "int* p = 3, ** q = 5;\n"
	result := builder.String()
	if result != expected {
		t.Fatalf("VariableDeclaration Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestParameterListWrite(t *testing.T) {
	param1 := NewParameterListItem(NewType("int", 0), "a")
	param2 := NewParameterListItem(NewType("float", 1), "b")

	paramList := NewParameterList(param1, param2)

	checkInterfaceCodeElement(param1)
	checkInterfaceCodeElement(param2)
	checkInterfaceCodeElement(paramList)

	builder, writer := makeTestWriter(testStyle1)
	err := paramList.Write(writer, 0)
	if err != nil {
		t.Fatalf("ParameterList Write failed: %s", err)
	}

	expected := "int a, float* b"
	result := builder.String()
	if result != expected {
		t.Fatalf("ParameterList Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
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

	builder, writer := makeTestWriter(testStyle1)
	err := f.Write(writer, 0)
	if err != nil {
		t.Fatalf("FunctionDeclaration Write failed: %s", err)
	}

	expected := strings.Join([]string{
		"int add(int a, int b) {",
		"}",
		"",
	}, "\n")
	result := builder.String()
	if result != expected {
		t.Fatalf("FunctionDeclaration Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
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

	builder, writer := makeTestWriter(testStyle1)
	err := f.Write(writer, 0)
	if err != nil {
		t.Fatalf("FunctionDeclaration Write failed: %s", err)
	}

	expected := strings.Join([]string{
		"int add(int a, int* b) {",
		"    return 42;",
		"}",
		"",
	}, "\n")
	result := builder.String()
	if result != expected {
		t.Fatalf("FunctionDeclaration Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
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

	builder, writer := makeTestWriter(testStyle2)
	err := f.Write(writer, 0)
	if err != nil {
		t.Fatalf("FunctionDeclaration Write failed: %s", err)
	}

	expected := strings.Join([]string{
		"int add(int a,int *b)",
		"{",
		"    return 42;",
		"}",
		"",
	}, "\n")
	result := builder.String()
	if result != expected {
		t.Fatalf("FunctionDeclaration Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}
