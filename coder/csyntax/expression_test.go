package csyntax

import (
	"testing"
)

type ExpressionTestCase struct {
	Result   Expression
	Expected string
}

type ExpressionTestCases []ExpressionTestCase

func (cases ExpressionTestCases) Run(t *testing.T, style *CodeStyle) {
	for _, testCase := range cases {
		checkInterfaceCodeElement(testCase.Result)
		checkInterfaceExpression(testCase.Result)

		checkOutputOnStyle(t, style, testCase.Expected, testCase.Result)
	}
}

func TestIdentifierWrite(t *testing.T) {
	id := NewIdentifier("myVar")

	checkInterfaceCodeElement(id)
	checkInterfaceExpression(id)

	expected := "myVar"
	checkOutputOnStyle(t, testStyle1, expected, id)
}

func TestIdentifierUnaryOperators(t *testing.T) {
	id := NewIdentifier("counter")

	ExpressionTestCases{
		{
			Result:   id.IncrPrefix(),
			Expected: "++counter",
		},
		{
			Result:   id.DecrPrefix(),
			Expected: "--counter",
		},
		{
			Result:   id.IncrPostfix(),
			Expected: "counter++",
		},
		{
			Result:   id.DecrPostfix(),
			Expected: "counter--",
		},
	}.Run(t, testStyle1)
}

func TestIdentifierInfixOperators(t *testing.T) {
	id1 := NewIdentifier("a")
	id2 := NewIdentifier("b")

	ExpressionTestCases{
		{
			Result:   id1.Add(id2),
			Expected: "a + b",
		},
		{
			Result:   id1.Sub(id2),
			Expected: "a - b",
		},
		{
			Result:   id1.Mul(id2),
			Expected: "a * b",
		},
		{
			Result:   id1.Div(id2),
			Expected: "a / b",
		},
	}.Run(t, testStyle1)
}

func TestInfixExpressionWrite(t *testing.T) {
	left := NewIdentifier("a")
	right := NewIdentifier("b")
	operator := OperatorAdd
	expr := NewInfixExpression(left, operator, right)

	checkInterfaceCodeElement(expr)
	checkInterfaceExpression(expr)

	expected := "a + b"
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

	expected := "(a + (b * c)) - d"
	checkOutputOnStyle(t, testStyle1, expected, fullExpr)
}

func TestPostfixExpressionWrite(t *testing.T) {
	operand := NewIdentifier("x")
	operator := OperatorIncrement
	expr := NewPostfixExpression(operand, operator)

	checkInterfaceCodeElement(expr)
	checkInterfaceExpression(expr)

	expected := "x++"
	checkOutputOnStyle(t, testStyle1, expected, expr)
}

func TestPostfixExpressionNested(t *testing.T) {
	operand := NewIdentifier("x")
	operator := OperatorIncrement

	expr1 := NewPostfixExpression(operand, operator)
	expected1 := "x++"
	checkInterfaceCodeElement(expr1)
	checkInterfaceExpression(expr1)
	checkOutputOnStyle(t, testStyle1, expected1, expr1)

	expr2 := NewPostfixExpression(expr1, operator)
	expected2 := "(x++)++"
	checkInterfaceCodeElement(expr2)
	checkInterfaceExpression(expr2)
	checkOutputOnStyle(t, testStyle1, expected2, expr2)
}

func TestUnaryExpressionWrite(t *testing.T) {
	operand := NewIdentifier("y")
	operator := OperatorNegative
	expr := NewUnaryExpression(operator, operand)

	checkInterfaceCodeElement(expr)
	checkInterfaceExpression(expr)

	expected := "-y"
	checkOutputOnStyle(t, testStyle1, expected, expr)
}

func TestUnaryExpressionNested(t *testing.T) {
	operand := NewIdentifier("y")
	operator := OperatorNegative

	expr1 := NewUnaryExpression(operator, operand)
	expected1 := "-y"
	checkInterfaceCodeElement(expr1)
	checkInterfaceExpression(expr1)
	checkOutputOnStyle(t, testStyle1, expected1, expr1)

	expr2 := NewUnaryExpression(operator, expr1)
	expected2 := "-(-y)"
	checkInterfaceCodeElement(expr2)
	checkInterfaceExpression(expr2)
	checkOutputOnStyle(t, testStyle1, expected2, expr2)
}
