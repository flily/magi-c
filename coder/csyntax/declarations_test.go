package csyntax

import (
	"testing"
)

func TestVariableDeclarationOneVariable(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("a", 0, NewIntegerLiteral(3))

	builder, writer := makeTestWriter(KRStyle)
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

func TestVariableDeclarationTwoVariables(t *testing.T) {
	stat := NewVariableDeclaration("int", nil)
	stat.Add("a", 0, NewIntegerLiteral(3))
	stat.Add("b", 0, NewIntegerLiteral(5))

	builder, writer := makeTestWriter(KRStyle)
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
