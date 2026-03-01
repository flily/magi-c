package csyntax

import (
	"testing"

	"strings"
)

func TestCodeBlockWrite(t *testing.T) {
	stat1 := NewAssignmentStatement("a", 0, NewIntegerLiteral(10))
	stat2 := NewAssignmentStatement("b", 1, NewIntegerLiteral(20))

	block := NewCodeBlock([]Statement{stat1, stat2})

	checkInterfaceCodeElement(block)
	checkInterfaceStatement(block)

	expected := strings.Join([]string{
		"    a = 10;\n",
		"    *b = 20;\n",
	}, "")
	checkOutputOnStyle(t, testStyle1, expected, block)
}

func TestEmptyCodeBlockWrite(t *testing.T) {
	l := NewEmptyLine()

	checkInterfaceCodeElement(l)
	checkInterfaceStatement(l)
	checkInterfaceDeclaration(l)

	expected := "\n"

	checkOutputOnStyle(t, testStyle1, expected, l)
}

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

func TestIfStatementWithoutElse1(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("a"), OperatorGreaterThan, NewIdentifier("b"))
	thenBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("a")),
	})

	ifStat := NewIfStatement(cond, thenBlock)

	checkInterfaceCodeElement(ifStat)
	checkInterfaceStatement(ifStat)

	expected := strings.Join([]string{
		"if ((a > b)) {",
		"    return a;",
		"}",
	}, "\n") + "\n"
	checkOutputOnStyle(t, testStyle1, expected, ifStat)
}

func TestIfStatementWithoutElse2(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("a"), OperatorGreaterThan, NewIdentifier("b"))
	thenBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("a")),
	})

	ifStat := NewIfStatement(cond, thenBlock)

	checkInterfaceCodeElement(ifStat)
	checkInterfaceStatement(ifStat)

	style := testStyle1.Clone()
	style.IfSpacing = false
	style.IfBraceOnNewLine = true
	style.IfBraceIndent = ""

	expected := strings.Join([]string{
		"if((a > b))",
		"{",
		"    return a;",
		"}",
	}, "\n") + "\n"
	checkOutputOnStyle(t, style, expected, ifStat)
}

func TestIfStatementWithIndent1(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("a"), OperatorGreaterThan, NewIdentifier("b"))
	thenBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("a")),
	})

	ifStat := NewIfStatement(cond, thenBlock)

	checkInterfaceCodeElement(ifStat)
	checkInterfaceStatement(ifStat)

	expected := strings.Join([]string{
		"    if ((a > b)) {",
		"        return a;",
		"    }",
	}, "\n") + "\n"
	checkOutputOnStyleWithIndentLevel(t, testStyle1, 1, expected, ifStat)
}

func TestIfStatementWittIndent2(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("a"), OperatorGreaterThan, NewIdentifier("b"))
	thenBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("a")),
	})

	ifStat := NewIfStatement(cond, thenBlock)

	checkInterfaceCodeElement(ifStat)
	checkInterfaceStatement(ifStat)

	style := testStyle1.Clone()
	style.IfSpacing = false
	style.IfBraceOnNewLine = true
	style.IfBraceIndent = ""

	expected := strings.Join([]string{
		"    if((a > b))",
		"    {",
		"        return a;",
		"    }",
	}, "\n") + "\n"
	checkOutputOnStyleWithIndentLevel(t, style, 1, expected, ifStat)
}
