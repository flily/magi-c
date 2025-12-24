package csyntax

import (
	"testing"
)

func TestVariableDeclarationOneVariableStyle1(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("a", 0, NewIntegerLiteral(3))

	var _ Statement = stat
	stat.statementNode()
	var _ Declaration = stat
	stat.declarationNode()

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

func TestVariableDeclarationTwoPointerVariableStyle2(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("p", 1, NewIntegerLiteral(3))
	stat.Add("q", 2, NewIntegerLiteral(5))

	builder, writer := makeTestWriter(testStyle2)
	err := stat.Write(writer, 0)
	if err != nil {
		t.Fatalf("VariableDeclaration Write failed: %s", err)
	}

	expected := "int *p = 3,**q = 5;\n"
	result := builder.String()
	if result != expected {
		t.Fatalf("VariableDeclaration Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}
