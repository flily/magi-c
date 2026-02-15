package csyntax

import (
	"testing"
)

func TestIdentifierWrite(t *testing.T) {
	id := NewIdentifier("myVar")

	checkInterfaceCodeElement(id)
	checkInterfaceExpression(id)

	expected := "myVar"
	checkOutputOnStyle(t, testStyle1, expected, id)
}

func TestInfixExpressionWrite(t *testing.T) {
	left := NewIdentifier("a")
	right := NewIdentifier("b")
	operator := OperatorAdd
	expr := NewInfixExpression(left, operator, right)

	checkInterfaceCodeElement(expr)
	checkInterfaceExpression(expr)

	expected := "(a + b)"
	checkOutputOnStyle(t, testStyle1, expected, expr)
}

func TestNestedInfixExpressionWrite(t *testing.T) {
	a := NewIdentifier("a")
	b := NewIdentifier("b")
	c := NewIdentifier("c")
	d := NewIdentifier("d")

	innerExpr := NewInfixExpression(b, OperatorMultiply, c)
	outerExpr := NewInfixExpression(a, OperatorAdd, innerExpr)
	fullExpr := NewInfixExpression(outerExpr, OperatorSubtract, d)

	checkInterfaceCodeElement(fullExpr)
	checkInterfaceExpression(fullExpr)

	expected := "((a + (b * c)) - d)"
	checkOutputOnStyle(t, testStyle1, expected, fullExpr)
}
