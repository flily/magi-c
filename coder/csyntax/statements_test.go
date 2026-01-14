package csyntax

import (
	"testing"
)

func TestAssignmentStatementOnNormalVariable(t *testing.T) {
	stat := NewAssignmentStatement("a", 0, NewIntegerLiteral(10))

	checkInterfaceCodeElement(stat)
	checkInterfaceStatement(stat)

	expected := "a = 10;\n"
	checkOutputOnStyle(t, testStyle1, expected, stat)
}

func TestAssignmentStatementOnPointerVariable(t *testing.T) {
	stat := NewAssignmentStatement("p", 1, NewIntegerLiteral(20))

	checkInterfaceCodeElement(stat)
	checkInterfaceStatement(stat)

	expected := "*p = 20;\n"
	checkOutputOnStyle(t, testStyle1, expected, stat)
}

func TestReturnStatementWithoutExpression(t *testing.T) {
	stat := NewReturnStatement(nil)

	checkInterfaceCodeElement(stat)
	checkInterfaceStatement(stat)

	expected := "return;\n"
	checkOutputOnStyle(t, testStyle1, expected, stat)
}

func TestReturnStatementWithSimpleIntegerLiteral(t *testing.T) {
	stat := NewReturnStatement(NewIntegerLiteral(42))

	checkInterfaceCodeElement(stat)
	checkInterfaceStatement(stat)

	expected := "return 42;\n"
	checkOutputOnStyle(t, testStyle1, expected, stat)
}
