package ast

import (
	"testing"

	"strings"
)

func TestSimpleExpressionList(t *testing.T) {
	text := "42 , 3.14 , PI"
	ctxList := generateTestWords(text)

	list := NewExpressionList()
	list.Add(NewIntegerLiteral(ctxList[0], 42), NewTerminalToken(ctxList[1], Comma))
	list.Add(NewFloatLiteral(ctxList[2], 3.14), NewTerminalToken(ctxList[3], Comma))
	list.Add(NewIdentifier(ctxList[4]), nil)

	expected := ASTBuildExpressionList(
		ASTBuildExpressionListItemWithComma(ASTBuildValue(42)),
		ASTBuildExpressionListItemWithComma(ASTBuildValue(3.14)),
		ASTBuildExpressionListItemWithoutComma(ASTBuildIdentifier("PI")),
	)

	if err := list.EqualTo(nil, expected); err != nil {
		t.Fatalf("ExpressionList not equal: %s", err)
	}
}

func TestExpressionListNotEqualInSize(t *testing.T) {
	text := "42 , 3.14 , PI"
	ctxList := generateTestWords(text)

	list := NewExpressionList()
	list.Add(NewIntegerLiteral(ctxList[0], 42), NewTerminalToken(ctxList[1], Comma))
	list.Add(NewFloatLiteral(ctxList[2], 3.14), NewTerminalToken(ctxList[3], Comma))
	list.Add(NewIdentifier(ctxList[4]), nil)

	expected := ASTBuildExpressionList(
		ASTBuildExpressionListItemWithComma(ASTBuildValue(42)),
		ASTBuildExpressionListItemWithComma(ASTBuildValue(3.14)),
	)
	message := strings.Join([]string{
		"   1:   42 , 3.14 , PI",
		"        ^^ ^ ^^^^ ^ ^^",
		"        wrong number of EXPRESSION LIST: expected 2, got 3",
	}, "\n")

	err := list.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("ExpressionList size mismatch not detected")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestExpressionListNotEqualInValueType(t *testing.T) {
	text := "42 , 3.14 , PI"
	ctxList := generateTestWords(text)

	list := NewExpressionList()
	list.Add(NewIntegerLiteral(ctxList[0], 42), NewTerminalToken(ctxList[1], Comma))
	list.Add(NewFloatLiteral(ctxList[2], 3.14), NewTerminalToken(ctxList[3], Comma))
	list.Add(NewIdentifier(ctxList[4]), nil)

	expected := ASTBuildExpressionList(
		ASTBuildExpressionListItemWithComma(ASTBuildValue(42)),
		ASTBuildExpressionListItemWithComma(ASTBuildValue(314)),
		ASTBuildExpressionListItemWithoutComma(ASTBuildIdentifier("PI")),
	)
	message := strings.Join([]string{
		"   1:   42 , 3.14 , PI",
		"             ^^^^",
		"             expect a *ast.IntegerLiteral",
	}, "\n")

	err := list.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("ExpressionList value type mismatch not detected")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestExpressionListNotEqualInValue(t *testing.T) {
	text := "42 , 3.14 , PI"
	ctxList := generateTestWords(text)

	list := NewExpressionList()
	list.Add(NewIntegerLiteral(ctxList[0], 42), NewTerminalToken(ctxList[1], Comma))
	list.Add(NewFloatLiteral(ctxList[2], 3.14), NewTerminalToken(ctxList[3], Comma))
	list.Add(NewIdentifier(ctxList[4]), nil)

	expected := ASTBuildExpressionList(
		ASTBuildExpressionListItemWithComma(ASTBuildValue(42)),
		ASTBuildExpressionListItemWithComma(ASTBuildValue(3.14)),
		ASTBuildExpressionListItemWithoutComma(ASTBuildIdentifier("pi")),
	)
	message := strings.Join([]string{
		"   1:   42 , 3.14 , PI",
		"                    ^^",
		"                    wrong identifier name, expect 'pi', got 'PI'",
	}, "\n")

	err := list.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("ExpressionList value mismatch not detected")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestExpressionListNotEqualInComma(t *testing.T) {
	text := "42 , 3.14 , PI"
	ctxList := generateTestWords(text)

	list := NewExpressionList()
	list.Add(NewIntegerLiteral(ctxList[0], 42), NewTerminalToken(ctxList[1], Comma))
	list.Add(NewFloatLiteral(ctxList[2], 3.14), NewTerminalToken(ctxList[3], Comma))
	list.Add(NewIdentifier(ctxList[4]), nil)

	expected := ASTBuildExpressionList(
		ASTBuildExpressionListItemWithComma(ASTBuildValue(42)),
		ASTBuildExpressionListItemWithComma(ASTBuildValue(3.14)),
		ASTBuildExpressionListItemWithComma(ASTBuildIdentifier("PI")),
	)
	message := strings.Join([]string{
		"   1:   42 , 3.14 , PI<EOF>",
		"                      ^^^^^",
		"                      expect *ast.TerminalToken",
	}, "\n")

	err := list.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("ExpressionList comma mismatch not detected")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestExpressionListItemNotEqualInType(t *testing.T) {
	text := "42 , 3.14 , PI"
	ctxList := generateTestWords(text)

	item := NewExpressionListItem(
		NewIntegerLiteral(ctxList[0], 42),
		NewTerminalToken(ctxList[1], Comma),
	)

	expected := ASTBuildValue(42)
	message := strings.Join([]string{
		"   1:   42 , 3.14 , PI",
		"        ^^ ^",
		"        expect a *ast.IntegerLiteral",
	}, "\n")

	err := item.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("ExpressionListItem type mismatch not detected")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s", err)
	}
}
func TestExpressionInfix(t *testing.T) {
	text := "a + b"
	ctxList := generateTestWords(text)

	expr := NewInfixExpression(
		NewIdentifier(ctxList[0]),
		NewTerminalToken(ctxList[1], Plus),
		NewIdentifier(ctxList[2]),
	)

	var _ Expression = expr
	expr.expressionNode()

	expected := ASTBuildInfixExpression(
		ASTBuildIdentifier("a"),
		Plus,
		ASTBuildIdentifier("b"),
	)

	if err := expr.EqualTo(nil, expected); err != nil {
		t.Fatalf("InfixExpression not equal: %s", err)
	}
}

func TestExpressionInfixNotEqualInType(t *testing.T) {
	text := "a + b"
	ctxList := generateTestWords(text)

	expr := NewInfixExpression(
		NewIdentifier(ctxList[0]),
		NewTerminalToken(ctxList[1], Plus),
		NewIdentifier(ctxList[2]),
	)

	expected := ASTBuildValue(42)
	message := strings.Join([]string{
		"   1:   a + b",
		"        ^ ^ ^",
		"        expect a *ast.IntegerLiteral",
	}, "\n")

	err := expr.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("InfixExpression type mismatch not detected")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestExpressionInfixNotEqualInLeft(t *testing.T) {
	text := "a + b"
	ctxList := generateTestWords(text)

	expr := NewInfixExpression(
		NewIdentifier(ctxList[0]),
		NewTerminalToken(ctxList[1], Plus),
		NewIdentifier(ctxList[2]),
	)

	expected := ASTBuildInfixExpression(
		ASTBuildIdentifier("x"),
		Plus,
		ASTBuildIdentifier("b"),
	)
	message := strings.Join([]string{
		"   1:   a + b",
		"        ^",
		"        wrong identifier name, expect 'x', got 'a'",
	}, "\n")

	err := expr.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("InfixExpression left operand mismatch not detected")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestExpressionInfixNotEqualInOperator(t *testing.T) {
	text := "a + b"
	ctxList := generateTestWords(text)

	expr := NewInfixExpression(
		NewIdentifier(ctxList[0]),
		NewTerminalToken(ctxList[1], Plus),
		NewIdentifier(ctxList[2]),
	)

	expected := ASTBuildInfixExpression(
		ASTBuildIdentifier("a"),
		Sub,
		ASTBuildIdentifier("b"),
	)
	message := strings.Join([]string{
		"   1:   a + b",
		"          ^",
		"          expect operator '-', got '+'",
	}, "\n")

	err := expr.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("InfixExpression operator mismatch not detected")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}

func TestExpressionInfixNotEqualInRight(t *testing.T) {
	text := "a + b"
	ctxList := generateTestWords(text)

	expr := NewInfixExpression(
		NewIdentifier(ctxList[0]),
		NewTerminalToken(ctxList[1], Plus),
		NewIdentifier(ctxList[2]),
	)

	expected := ASTBuildInfixExpression(
		ASTBuildIdentifier("a"),
		Plus,
		ASTBuildIdentifier("y"),
	)
	message := strings.Join([]string{
		"   1:   a + b",
		"            ^",
		"            wrong identifier name, expect 'y', got 'b'",
	}, "\n")

	err := expr.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("InfixExpression right operand mismatch not detected")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\n%s\nexpect\n%s", err, message)
	}
}
