package ast

import (
	"testing"

	"strings"

	"github.com/flily/magi-c/context"
)

func makeFunctionDeclarationTokens1(ctxList []*context.Context) *FunctionDeclaration {
	funcToken := NewTerminalToken(ctxList[0], Function)
	funcDecl := NewFunctionDeclaration(funcToken)

	funcDecl.Name = NewIdentifier(ctxList[1])
	funcDecl.LParenArgs = NewTerminalToken(ctxList[2], LeftParen)

	args := NewArgumentList()

	argAName := NewIdentifier(ctxList[3])
	argAType := NewSimpleType()
	argAType.Identifier = NewIdentifier(ctxList[4])
	argAComma := NewTerminalToken(ctxList[5], Comma)
	args.Add(argAName, argAType, argAComma)

	argBName := NewIdentifier(ctxList[6])
	argBType := NewSimpleType()
	argBType.Identifier = NewIdentifier(ctxList[7])
	args.Add(argBName, argBType, nil)

	funcDecl.Arguments = args
	funcDecl.RParenArgs = NewTerminalToken(ctxList[8], RightParen)

	funcDecl.LParenReturnTypes = NewTerminalToken(ctxList[9], LeftParen)

	returnType := NewSimpleType()
	returnType.Identifier = NewIdentifier(ctxList[10])

	returnTypes := NewTypeList()
	returnTypes.Add(returnType, nil)
	funcDecl.ReturnTypes = returnTypes

	funcDecl.RParenReturnTypes = NewTerminalToken(ctxList[11], RightParen)

	funcDecl.LBrace = NewTerminalToken(ctxList[12], LeftBrace)

	returnToken := NewTerminalToken(ctxList[13], Return)
	returnStmt := NewReturnStatement(returnToken)

	returnValues := NewExpressionList()

	returnValueA := NewIdentifier(ctxList[14])
	returnValues.Add(returnValueA, NewTerminalToken(ctxList[15], Comma))

	returnValueB := NewIdentifier(ctxList[16])
	returnValues.Add(returnValueB, nil)

	returnStmt.Value = returnValues

	funcDecl.Statements = []Statement{returnStmt}

	funcDecl.RBrace = NewTerminalToken(ctxList[17], RightBrace)

	return funcDecl
}

func TestFunctionDeclaration(t *testing.T) {
	text := strings.Join([]string{
		//  0   1 2 3   4 5 6   7 8 9 10 11 12
		"func add ( a int , b int ) ( int ) {",
		//      13 14 15 16
		"    return a + b",
		// 17
		"}",
	}, "\n")
	ctxList := generateTestWords(text)

	funcDecl := makeFunctionDeclarationTokens1(ctxList)

	var _ Declaration = funcDecl
	funcDecl.declarationNode()

	expected := ASTBuildFunction(
		"add",
		ASTBuildArgumentList(
			ASTBuildArgumentWithComma("a", ASTBuildSimpleType("int")),
			ASTBuildArgumentWithoutComma("b", ASTBuildSimpleType("int")),
		),
		ASTBuildTypeList(
			ASTBuildTypeListItemWithoutComma(ASTBuildSimpleType("int")),
		),
		[]Statement{
			ASTBuildReturnStatement(
				ASTBuildExpressionList(
					ASTBuildExpressionListItemWithComma(
						ASTBuildIdentifier("a"),
					),
					ASTBuildExpressionListItemWithoutComma(
						ASTBuildIdentifier("b"),
					),
				),
			),
		},
	)

	if err := funcDecl.EqualTo(nil, expected); err != nil {
		t.Errorf("FunctionDeclaration not equal:\n%s", err)
	}
}

func TestFunctionDeclarationNotEqual(t *testing.T) {
	text := strings.Join([]string{
		//  0   1 2 3   4 5 6   7 8 9 10 11 12
		"func add ( a int , b int ) ( int ) {",
		//      13 14 15 16
		"    return b + a",
		// 17
		"}",
	}, "\n")
	ctxList := generateTestWords(text)

	funcDecl := makeFunctionDeclarationTokens1(ctxList)

	var _ Declaration = funcDecl
	funcDecl.declarationNode()

	expected := ASTBuildFunction(
		"add",
		ASTBuildArgumentList(
			ASTBuildArgumentWithComma("a", ASTBuildSimpleType("int")),
			ASTBuildArgumentWithoutComma("b", ASTBuildSimpleType("int")),
		),
		ASTBuildTypeList(
			ASTBuildTypeListItemWithoutComma(ASTBuildSimpleType("int")),
		),
		[]Statement{
			ASTBuildReturnStatement(
				ASTBuildExpressionList(
					ASTBuildExpressionListItemWithComma(
						ASTBuildIdentifier("a"),
					),
					ASTBuildExpressionListItemWithoutComma(
						ASTBuildIdentifier("b"),
					),
				),
			),
		},
	)

	message := strings.Join([]string{
		"   2:       return b + a",
		"                   ^",
		"                   wrong identifier name, expect a, got b",
	}, "\n")
	err := funcDecl.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("FunctionDeclaration expected not equal, but equal")
	}

	if err.Error() != message {
		t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
	}
}
