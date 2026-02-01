package ast

import (
	"testing"

	"strings"
)

func TestReturnStatement(t *testing.T) {
	text := "return 0"
	ctxList := generateTestWords(text)

	returnToken := NewTerminalToken(ctxList[0], Return)
	returnStmt := NewReturnStatement(returnToken)

	returnValues := NewExpressionList()
	returnValues.Add(NewIntegerLiteral(ctxList[1], 0), nil)
	returnStmt.Value = returnValues

	var _ Statement = returnStmt
	returnStmt.statementNode()

	expected := ASTBuildReturnStatement(
		ASTBuildExpressionList(
			ASTBuildExpressionListItemWithoutComma(ASTBuildValue(0)),
		),
	)

	if err := returnStmt.EqualTo(nil, expected); err != nil {
		t.Errorf("ReturnStatement not equal:\n%s", err)
	}
}

func TestReturnStatementNotEqualInReturnValue(t *testing.T) {
	text := "return 0"
	ctxList := generateTestWords(text)

	returnToken := NewTerminalToken(ctxList[0], Return)
	returnStmt := NewReturnStatement(returnToken)

	returnValues := NewExpressionList()
	returnValues.Add(NewIntegerLiteral(ctxList[1], 0), nil)
	returnStmt.Value = returnValues

	expected := ASTBuildReturnStatement(
		ASTBuildExpressionList(
			ASTBuildExpressionListItemWithoutComma(ASTBuildValue(1)),
		),
	)
	message := strings.Join([]string{
		"test.txt:1:8: error: wrong integer value, expect 1, got 0",
		"    1 | return 0",
		"      |        ^",
		"      |        1",
	}, "\n")

	err := returnStmt.EqualTo(nil, expected)
	if err == nil {
		t.Errorf("ReturnStatement expected not equal, but equal")
	}

	if err.Error() != message {
		t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
	}
}

func TestReturnStatementNotEqual(t *testing.T) {
	text := "return 0"
	ctxList := generateTestWords(text)

	returnToken := NewTerminalToken(ctxList[0], Return)
	returnStmt := NewReturnStatement(returnToken)

	returnValues := NewExpressionList()
	returnValues.Add(NewIntegerLiteral(ctxList[1], 0), nil)
	returnStmt.Value = returnValues

	expected := ASTBuildValue(0)
	message := strings.Join([]string{
		"test.txt:1:1: error: expect a *ast.IntegerLiteral, got a *ast.ReturnStatement",
		"    1 | return 0",
		"      | ^^^^^^ ^",
		"      | *ast.IntegerLiteral",
	}, "\n")

	err := returnStmt.EqualTo(nil, expected)
	if err == nil {
		t.Errorf("ReturnStatement expected not equal, but equal")
	}

	if err.Error() != message {
		t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
	}
}
