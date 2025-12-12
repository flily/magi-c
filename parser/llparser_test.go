package parser

import (
	"testing"

	"strings"

	"github.com/flily/magi-c/ast"
)

func TestLLParserSimpleStatement(t *testing.T) {
	newCorrectCodeTestCase(
		strings.Join([]string{
			"#include <stdio.h>",
		}, "\n"),
		ast.ASTBuildDocument(
			ast.ASTBuildIncludeAngle("stdio.h"),
		),
	).Run(t)
}

func TestLLParserSimplestProgram(t *testing.T) {
	newCorrectCodeTestCase(
		strings.Join([]string{
			"fun main() (int) {",
			"    return 0",
			"}",
		}, "\n"),
		ast.ASTBuildDocument(
			ast.ASTBuildFunction(
				"main",
				nil,
				ast.ASTBuildTypeList(
					ast.ASTBuildTypeListItemWithoutComma("int"),
				),
				[]ast.Statement{
					ast.ASTBuildReturnStatement(
						ast.ASTBuildExpressionList(
							ast.ASTBuildExpressionListItemWithoutComma(
								ast.ASTBuildValue(0),
							),
						),
					),
				},
			),
		),
	).Run(t)
}

func TestLLParserFunctionWithArguments(t *testing.T) {
	newCorrectCodeTestCase(
		strings.Join([]string{
			"fun add(a int, b int) (int, int) {",
			"    return 0, 0",
			"}",
		}, "\n"),

		ast.ASTBuildDocument(
			ast.ASTBuildFunction(
				"add",
				ast.ASTBuildArgumentList(
					ast.ASTBuildArgumentWithComma("a", "int"),
					ast.ASTBuildArgumentWithoutComma("b", "int"),
				),
				ast.ASTBuildTypeList(
					ast.ASTBuildTypeListItemWithComma("int"),
					ast.ASTBuildTypeListItemWithoutComma("int"),
				),
				[]ast.Statement{
					ast.ASTBuildReturnStatement(
						ast.ASTBuildExpressionList(
							ast.ASTBuildExpressionListItemWithComma(
								ast.ASTBuildValue(0),
							),
							ast.ASTBuildExpressionListItemWithoutComma(
								ast.ASTBuildValue(0),
							),
						),
					),
				},
			),
		),
	).Run(t)
}

func TestLLParserReturnWithExpressionList(t *testing.T) {
	newCorrectCodeTestCase(
		strings.Join([]string{
			"fun foo() (int, int, int) {",
			"    return 1, 2, 3",
			"}",
		}, "\n"),
		ast.ASTBuildDocument(
			ast.ASTBuildFunction(
				"foo",
				nil,
				ast.ASTBuildTypeList(
					ast.ASTBuildTypeListItemWithComma("int"),
					ast.ASTBuildTypeListItemWithComma("int"),
					ast.ASTBuildTypeListItemWithoutComma("int"),
				),
				[]ast.Statement{
					ast.ASTBuildReturnStatement(
						ast.ASTBuildExpressionList(
							ast.ASTBuildExpressionListItemWithComma(
								ast.ASTBuildValue(1),
							),
							ast.ASTBuildExpressionListItemWithComma(
								ast.ASTBuildValue(2),
							),
							ast.ASTBuildExpressionListItemWithoutComma(
								ast.ASTBuildValue(3),
							),
						),
					),
				},
			),
		),
	).Run(t)
}

func TestLLParserReturnWithExpressionArithmetic1(t *testing.T) {
	newCorrectCodeTestCase(
		strings.Join([]string{
			"fun add() (int) {",
			"    return 1 + 2",
			"}",
		}, "\n"),
		ast.ASTBuildDocument(
			ast.ASTBuildFunction(
				"add",
				nil,
				ast.ASTBuildTypeList(
					ast.ASTBuildTypeListItemWithoutComma("int"),
				),
				[]ast.Statement{
					ast.ASTBuildReturnStatement(
						ast.ASTBuildExpressionList(
							ast.ASTBuildExpressionListItemWithoutComma(
								ast.ASTBuildInfixExpression(
									ast.ASTBuildValue(1),
									ast.Plus,
									ast.ASTBuildValue(2),
								),
							),
						),
					),
				},
			),
		),
	).Run(t)
}

func TestLLParserReturnWithExpressionArithmetic2(t *testing.T) {
	newCorrectCodeTestCase(
		strings.Join([]string{
			"fun add(a int, b int) (int) {",
			"    return a + b",
			"}",
		}, "\n"),

		ast.ASTBuildDocument(
			ast.ASTBuildFunction(
				"add",
				ast.ASTBuildArgumentList(
					ast.ASTBuildArgumentWithComma("a", "int"),
					ast.ASTBuildArgumentWithoutComma("b", "int"),
				),
				ast.ASTBuildTypeList(
					ast.ASTBuildTypeListItemWithoutComma("int"),
				),
				[]ast.Statement{
					ast.ASTBuildReturnStatement(
						ast.ASTBuildExpressionList(
							ast.ASTBuildExpressionListItemWithoutComma(
								ast.ASTBuildInfixExpression(
									ast.ASTBuildIdentifier("a"),
									ast.Plus,
									ast.ASTBuildIdentifier("b"),
								),
							),
						),
					),
				},
			),
		),
	).Run(t)
}

func TestLLParserReturnWithExpressionArithmetic3(t *testing.T) {
	newCorrectCodeTestCase(
		strings.Join([]string{
			"fun add(a int, b int) (int) {",
			"    return a + b + 3",
			"}",
		}, "\n"),
		ast.ASTBuildDocument(
			ast.ASTBuildFunction(
				"add",
				ast.ASTBuildArgumentList(
					ast.ASTBuildArgumentWithComma("a", "int"),
					ast.ASTBuildArgumentWithoutComma("b", "int"),
				),
				ast.ASTBuildTypeList(
					ast.ASTBuildTypeListItemWithoutComma("int"),
				),
				[]ast.Statement{
					ast.ASTBuildReturnStatement(
						ast.ASTBuildExpressionList(
							ast.ASTBuildExpressionListItemWithoutComma(
								ast.ASTBuildInfixExpression(
									ast.ASTBuildInfixExpression(
										ast.ASTBuildIdentifier("a"),
										ast.Plus,
										ast.ASTBuildIdentifier("b"),
									),
									ast.Plus,
									ast.ASTBuildValue(3),
								),
							),
						),
					),
				},
			),
		),
	).Run(t)
}
