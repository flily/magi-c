package csyntax

import (
	"testing"
)

func TestAssignmentStatementOnNormalVariable(t *testing.T) {
	stat := NewAssignmentStatement("a", 0, NewIntegerLiteral(10))

	checkInterfaceCodeElement(stat)
	checkInterfaceStatement(stat)

	builder, writer := makeTestWriter(testStyle1)
	err := stat.Write(writer, 0)
	if err != nil {
		t.Fatalf("AssignmentStatement Write failed: %s", err)
	}

	expected := "a = 10;\n"
	result := builder.String()
	if result != expected {
		t.Fatalf("AssignmentStatement Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestAssignmentStatementOnPointerVariable(t *testing.T) {
	stat := NewAssignmentStatement("p", 1, NewIntegerLiteral(20))
	var _ Statement = stat
	stat.statementNode()

	builder, writer := makeTestWriter(testStyle1)
	err := stat.Write(writer, 0)
	if err != nil {
		t.Fatalf("AssignmentStatement Write failed: %s", err)
	}

	expected := "*p = 20;\n"
	result := builder.String()
	if result != expected {
		t.Fatalf("AssignmentStatement Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestReturnStatementWithoutExpression(t *testing.T) {
	stat := NewReturnStatement(nil)

	checkInterfaceCodeElement(stat)
	checkInterfaceStatement(stat)

	builder, writer := makeTestWriter(testStyle1)
	err := stat.Write(writer, 0)
	if err != nil {
		t.Fatalf("ReturnStatement Write failed: %s", err)
	}

	expected := "return;\n"
	result := builder.String()
	if result != expected {
		t.Fatalf("ReturnStatement Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestReturnStatementWithSimpleIntegerLiteral(t *testing.T) {
	stat := NewReturnStatement(NewIntegerLiteral(42))

	builder, writer := makeTestWriter(testStyle1)
	err := stat.Write(writer, 0)
	if err != nil {
		t.Fatalf("ReturnStatement Write failed: %s", err)
	}

	expected := "return 42;\n"
	result := builder.String()
	if result != expected {
		t.Fatalf("ReturnStatement Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}
