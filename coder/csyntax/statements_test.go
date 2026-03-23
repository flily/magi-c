package csyntax

import (
	"testing"

	"strings"
)

func TestCodeBlockWrite(t *testing.T) {
	stmt1 := NewAssignmentStatement("a", 0, NewIntegerLiteral(10))
	stmt2 := NewAssignmentStatement("b", 1, NewIntegerLiteral(20))

	block := NewCodeBlock([]Statement{stmt1, stmt2})

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

	expected := "\n"

	checkOutputOnStyle(t, testStyle1, expected, l)
}

func TestCodeSegmentWrite(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("a", 0, NewIntegerLiteral(10))
	decl.Add("b", 0, NewIntegerLiteral(20))
	stmt1 := NewDeclarationStatement(decl)

	stmt2 := NewIfStatement(
		NewInfixExpression(NewIdentifier("a"), OperatorLessThan, NewIdentifier("b")),
		NewCodeBlock([]Statement{
			NewAssignmentStatement("a", 0, NewIntegerLiteral(0)),
		}),
	)

	segment := NewCodeSegment([]Statement{
		stmt1,
		stmt2,
	})

	checkInterfaceCodeElement(segment)
	checkInterfaceStatement(segment)

	expected := strings.Join([]string{
		"int a = 10, b = 20;",
		"if (a < b) {",
		"    a = 0;",
		"}",
		"",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle1, expected, segment)
}

func TestDeclarationStatmentOneVariableStyle1(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("a", 0, NewIntegerLiteral(3))

	stmt := NewDeclarationStatement(decl)
	checkInterfaceCodeElement(stmt)
	checkInterfaceStatement(stmt)

	expected := "int a = 3;\n"
	checkOutputOnStyle(t, testStyle1, expected, stmt)
}

func TestDeclarationStatmentOneVariableStyle2(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("a", 0, NewIntegerLiteral(3))

	stmt := NewDeclarationStatement(decl)
	checkInterfaceCodeElement(stmt)
	checkInterfaceStatement(stmt)

	expected := "int a = 3;\n"
	checkOutputOnStyle(t, testStyle2, expected, stmt)
}

func TestDeclarationStatmentTwoVariablesStyle1(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("a", 0, NewIntegerLiteral(3))
	decl.Add("b", 0, NewIntegerLiteral(5))

	stmt := NewDeclarationStatement(decl)
	checkInterfaceCodeElement(stmt)
	checkInterfaceStatement(stmt)

	expected := "int a = 3, b = 5;\n"
	checkOutputOnStyle(t, testStyle1, expected, stmt)
}

func TestDeclarationStatmentTwoVariablesStyle2(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("a", 0, NewIntegerLiteral(3))
	decl.Add("b", 0, NewIntegerLiteral(5))
	stmt := NewDeclarationStatement(decl)

	expected := "int a = 3,b = 5;\n"
	checkOutputOnStyle(t, testStyle2, expected, stmt)
}

func TestDeclarationStatmentOnePointerVariableStyle1(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("p", 1, NewIntegerLiteral(3))

	stmt := NewDeclarationStatement(decl)

	expected := "int* p = 3;\n"
	checkOutputOnStyle(t, testStyle1, expected, stmt)
}

func TestDeclarationStatmentOnePointerVariableStyle2(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("p", 1, NewIntegerLiteral(3))

	stmt := NewDeclarationStatement(decl)
	checkInterfaceCodeElement(stmt)
	checkInterfaceStatement(stmt)

	expected := "int *p = 3;\n"
	checkOutputOnStyle(t, testStyle2, expected, stmt)
}

func TestDeclarationStatmentTwoPointerVariableStyle1(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("p", 1, NewIntegerLiteral(3))
	decl.Add("q", 2, NewIntegerLiteral(5))

	stmt := NewDeclarationStatement(decl)
	checkInterfaceCodeElement(stmt)
	checkInterfaceStatement(stmt)

	expected := "int* p = 3, ** q = 5;\n"
	checkOutputOnStyle(t, testStyle1, expected, stmt)
}

func TestDeclarationStatmentTwoPointerVariableStyle2(t *testing.T) {
	decl := NewVariableDeclaration("int", nil)
	decl.Add("p", 1, NewIntegerLiteral(3))
	decl.Add("q", 2, NewIntegerLiteral(5))

	stmt := NewDeclarationStatement(decl)
	checkInterfaceCodeElement(stmt)
	checkInterfaceStatement(stmt)

	expected := "int *p = 3, **q = 5;\n"
	checkOutputOnStyle(t, testStyle2, expected, stmt)
}

func TestAssignmentStatementOnNormalVariableStyle1(t *testing.T) {
	stmt := NewAssignmentStatement("a", 0, NewIntegerLiteral(10))

	checkInterfaceCodeElement(stmt)
	checkInterfaceStatement(stmt)

	expected := "a = 10;\n"
	checkOutputOnStyle(t, testStyle1, expected, stmt)
}

func TestAssignmentStatementOnNormalVariableStyle2(t *testing.T) {
	stmt := NewAssignmentStatement("a", 0, NewIntegerLiteral(10))

	checkInterfaceCodeElement(stmt)
	checkInterfaceStatement(stmt)

	expected := "a = 10;\n"
	checkOutputOnStyle(t, testStyle2, expected, stmt)
}

func TestAssignmentStatementOnPointerVariableStyle1(t *testing.T) {
	stmt := NewAssignmentStatement("p", 1, NewIntegerLiteral(20))

	checkInterfaceCodeElement(stmt)
	checkInterfaceStatement(stmt)

	expected := "*p = 20;\n"
	checkOutputOnStyle(t, testStyle1, expected, stmt)
}

func TestAssignmentStatementOnPointerVariableStyle2(t *testing.T) {
	stmt := NewAssignmentStatement("p", 1, NewIntegerLiteral(20))

	checkInterfaceCodeElement(stmt)
	checkInterfaceStatement(stmt)

	expected := "* p = 20;\n"
	checkOutputOnStyle(t, testStyle2, expected, stmt)
}

func TestReturnStatementWithoutExpression(t *testing.T) {
	stmt := NewReturnStatement(nil)

	checkInterfaceCodeElement(stmt)
	checkInterfaceStatement(stmt)

	expected := "return;\n"
	checkOutputOnStyle(t, testStyle1, expected, stmt)
}

func TestReturnStatementWithSimpleIntegerLiteral(t *testing.T) {
	stmt := NewReturnStatement(NewIntegerLiteral(42))

	checkInterfaceCodeElement(stmt)
	checkInterfaceStatement(stmt)

	expected := "return 42;\n"
	checkOutputOnStyle(t, testStyle1, expected, stmt)
}

func TestIfStatementWithoutElseStyle1(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("a"), OperatorGreaterThan, NewIdentifier("b"))
	thenBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("a")),
	})

	ifStmt := NewIfStatement(cond, thenBlock)

	checkInterfaceCodeElement(ifStmt)
	checkInterfaceStatement(ifStmt)

	expected := strings.Join([]string{
		"if (a > b) {",
		"    return a;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle1, expected, ifStmt)
}

func TestIfStatementWithElseStyle1(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("a"), OperatorGreaterThan, NewIdentifier("b"))
	thenBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("a")),
	})
	elseBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("b")),
	})

	ifStat := NewIfElseStatement(cond, thenBlock, elseBlock)

	checkInterfaceCodeElement(ifStat)
	checkInterfaceStatement(ifStat)

	expected := strings.Join([]string{
		"if (a > b) {",
		"    return a;",
		"} else {",
		"    return b;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle1, expected, ifStat)
}

func TestIfStatementWithoutElseStyle2(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("a"), OperatorGreaterThan, NewIdentifier("b"))
	thenBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("a")),
	})

	ifStmt := NewIfStatement(cond, thenBlock)

	checkInterfaceCodeElement(ifStmt)
	checkInterfaceStatement(ifStmt)

	style := testStyle1.Clone()
	style.IfSpacing = false
	style.IfBraceOnNewLine = true
	style.IfBraceIndent = ""

	expected := strings.Join([]string{
		"if(a > b)",
		"{",
		"    return a;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, style, expected, ifStmt)
}

func TestIfStatementWithElseStyle2(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("a"), OperatorGreaterThan, NewIdentifier("b"))
	thenBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("a")),
	})
	elseBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("b")),
	})

	ifStmt := NewIfElseStatement(cond, thenBlock, elseBlock)

	checkInterfaceCodeElement(ifStmt)
	checkInterfaceStatement(ifStmt)

	expected := strings.Join([]string{
		"if(a > b)",
		"{",
		"    return a;",
		"}",
		"else",
		"{",
		"    return b;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle2, expected, ifStmt)
}

func TestIfStatementWithIndentStyle1(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("a"), OperatorGreaterThan, NewIdentifier("b"))
	thenBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("a")),
	})

	ifStmt := NewIfStatement(cond, thenBlock)

	checkInterfaceCodeElement(ifStmt)
	checkInterfaceStatement(ifStmt)

	expected := strings.Join([]string{
		"    if (a > b) {",
		"        return a;",
		"    }",
		"",
	}, "\n")
	level := NewLevel(1, 0)
	checkOutputOnStyleWithIndentLevel(t, testStyle1, level, expected, ifStmt)
}

func TestIfElseStatementWithIndentStyle1(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("a"), OperatorGreaterThan, NewIdentifier("b"))
	thenBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("a")),
	})
	elseBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("b")),
	})

	ifStmt := NewIfElseStatement(cond, thenBlock, elseBlock)

	checkInterfaceCodeElement(ifStmt)
	checkInterfaceStatement(ifStmt)

	expected := strings.Join([]string{
		"    if (a > b) {",
		"        return a;",
		"    } else {",
		"        return b;",
		"    }",
		"",
	}, "\n")
	level := NewLevel(1, 0)
	checkOutputOnStyleWithIndentLevel(t, testStyle1, level, expected, ifStmt)
}

func TestIfStatementWithIndentStyle2(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("a"), OperatorGreaterThan, NewIdentifier("b"))
	thenBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("a")),
	})

	ifStmt := NewIfStatement(cond, thenBlock)

	checkInterfaceCodeElement(ifStmt)
	checkInterfaceStatement(ifStmt)

	style := testStyle1.Clone()
	style.IfSpacing = false
	style.IfBraceOnNewLine = true
	style.IfBraceIndent = ""

	expected := strings.Join([]string{
		"    if(a > b)",
		"    {",
		"        return a;",
		"    }",
		"",
	}, "\n")
	level := NewLevel(1, 0)
	checkOutputOnStyleWithIndentLevel(t, style, level, expected, ifStmt)
}

func TestIfElseStatementWithIndentStyle2(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("a"), OperatorGreaterThan, NewIdentifier("b"))
	thenBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("a")),
	})
	elseBlock := NewCodeBlock([]Statement{
		NewReturnStatement(NewIdentifier("b")),
	})

	ifStmt := NewIfElseStatement(cond, thenBlock, elseBlock)

	checkInterfaceCodeElement(ifStmt)
	checkInterfaceStatement(ifStmt)

	style := testStyle1.Clone()
	style.IfSpacing = false
	style.IfBraceOnNewLine = true
	style.IfBraceIndent = ""

	expected := strings.Join([]string{
		"    if(a > b)",
		"    {",
		"        return a;",
		"    }",
		"    else",
		"    {",
		"        return b;",
		"    }",
		"",
	}, "\n")
	level := NewLevel(1, 0)
	checkOutputOnStyleWithIndentLevel(t, style, level, expected, ifStmt)
}

func TestIfElseChainStatementStyle2(t *testing.T) {
	cond1 := NewInfixExpression(NewIdentifier("a"), OperatorEqual, NewIntegerLiteral(1))
	thenBlock1 := NewCodeBlock([]Statement{
		NewAssignmentStatement("r", 0, NewIntegerLiteral(1)),
	})
	cond2 := NewInfixExpression(NewIdentifier("a"), OperatorEqual, NewIntegerLiteral(2))
	thenBlock2 := NewCodeBlock([]Statement{
		NewAssignmentStatement("r", 0, NewIntegerLiteral(2)),
	})
	cond3 := NewInfixExpression(NewIdentifier("a"), OperatorEqual, NewIntegerLiteral(3))
	thenBlock3 := NewCodeBlock([]Statement{
		NewAssignmentStatement("r", 0, NewIntegerLiteral(3)),
	})
	cond4 := NewInfixExpression(NewIdentifier("a"), OperatorEqual, NewIntegerLiteral(4))
	thenBlock4 := NewCodeBlock([]Statement{
		NewAssignmentStatement("r", 0, NewIntegerLiteral(4)),
	})
	elseBlock := NewCodeBlock([]Statement{
		NewAssignmentStatement("r", 0, NewIntegerLiteral(0)),
	})

	ifStmt := NewIfElseChainStatement(
		[]*IfStatement{
			NewIfStatement(cond1, thenBlock1),
			NewIfStatement(cond2, thenBlock2),
			NewIfStatement(cond3, thenBlock3),
			NewIfStatement(cond4, thenBlock4),
		},
		elseBlock,
	)

	checkInterfaceCodeElement(ifStmt)
	checkInterfaceStatement(ifStmt)

	expected := strings.Join([]string{
		"if (a == 1) {",
		"    r = 1;",
		"} else if (a == 2) {",
		"    r = 2;",
		"} else if (a == 3) {",
		"    r = 3;",
		"} else if (a == 4) {",
		"    r = 4;",
		"} else {",
		"    r = 0;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle1, expected, ifStmt)
}

func TestWhileStatementStyle1(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("i"), OperatorLessThan, NewIntegerLiteral(10))
	body := NewCodeBlock([]Statement{
		NewAssignmentStatement("i", 0, NewInfixExpression(NewIdentifier("i"), OperatorAdd, NewIntegerLiteral(1))),
	})

	whileStmt := NewWhileStatement(cond, body)

	checkInterfaceCodeElement(whileStmt)
	checkInterfaceStatement(whileStmt)

	expected := strings.Join([]string{
		"while (i < 10) {",
		"    i = i + 1;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle1, expected, whileStmt)
}

func TestWhileStatementStyle2(t *testing.T) {
	cond := NewInfixExpression(NewIdentifier("i"), OperatorLessThan, NewIntegerLiteral(10))
	body := NewCodeBlock([]Statement{
		NewAssignmentStatement("i", 0, NewInfixExpression(NewIdentifier("i"), OperatorAdd, NewIntegerLiteral(1))),
	})

	whileStmt := NewWhileStatement(cond, body)

	checkInterfaceCodeElement(whileStmt)
	checkInterfaceStatement(whileStmt)

	expected := strings.Join([]string{
		"while(i < 10)",
		"{",
		"    i = i + 1;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle2, expected, whileStmt)
}

func TestDoWhileStatementStyle1(t *testing.T) {
	body := NewCodeBlock([]Statement{
		NewAssignmentStatement("i", 0, NewInfixExpression(NewIdentifier("i"), OperatorAdd, NewIntegerLiteral(1))),
	})
	cond := NewInfixExpression(NewIdentifier("i"), OperatorLessThan, NewIntegerLiteral(10))

	doWhileStmt := NewDoWhileStatement(body, cond)

	checkInterfaceCodeElement(doWhileStmt)
	checkInterfaceStatement(doWhileStmt)

	expected := strings.Join([]string{
		"do {",
		"    i = i + 1;",
		"} while (i < 10);",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle1, expected, doWhileStmt)
}

func TestDoWhileStatementStyle2(t *testing.T) {
	body := NewCodeBlock([]Statement{
		NewAssignmentStatement("i", 0, NewInfixExpression(NewIdentifier("i"), OperatorAdd, NewIntegerLiteral(1))),
	})
	cond := NewInfixExpression(NewIdentifier("i"), OperatorLessThan, NewIntegerLiteral(10))

	doWhileStmt := NewDoWhileStatement(body, cond)

	checkInterfaceCodeElement(doWhileStmt)
	checkInterfaceStatement(doWhileStmt)

	expected := strings.Join([]string{
		"do",
		"{",
		"    i = i + 1;",
		"} while(i < 10);",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle2, expected, doWhileStmt)
}

func TestForStatementWithVariableDeclarationStyle1(t *testing.T) {
	initor := NewVariableDeclaration("int", nil)
	initor.Add("i", 0, NewIntegerLiteral(0))
	cond := NewInfixExpression(NewIdentifier("i"), OperatorLessThan, NewIntegerLiteral(10))
	update := NewIdentifier("i").IncrPostfix()
	body := NewCodeBlock([]Statement{
		NewAssignmentStatement("sum", 0, NewInfixExpression(NewIdentifier("sum"), OperatorAdd, NewIdentifier("i"))),
	})

	forStmt := NewForStatement(initor, cond, update, body)

	checkInterfaceCodeElement(forStmt)
	checkInterfaceStatement(forStmt)

	expected := strings.Join([]string{
		"for (int i = 0; i < 10; i++) {",
		"    sum = sum + i;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle1, expected, forStmt)
}

func TestForStatementWithVariableDeclarationStyle2(t *testing.T) {
	initor := NewVariableDeclaration("int", nil)
	initor.Add("i", 0, NewIntegerLiteral(0))
	cond := NewInfixExpression(NewIdentifier("i"), OperatorLessThan, NewIntegerLiteral(10))
	update := NewIdentifier("i").IncrPostfix()
	body := NewCodeBlock([]Statement{
		NewAssignmentStatement("sum", 0, NewInfixExpression(NewIdentifier("sum"), OperatorAdd, NewIdentifier("i"))),
	})

	forStmt := NewForStatement(initor, cond, update, body)

	checkInterfaceCodeElement(forStmt)
	checkInterfaceStatement(forStmt)

	expected := strings.Join([]string{
		"for(int i = 0; i < 10; i++)",
		"{",
		"    sum = sum + i;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle2, expected, forStmt)
}

func TestForStatementWithoutVariableDeclarationStyle1(t *testing.T) {
	initor := NewAssignmentExpression("i", 0, NewIntegerLiteral(0))
	cond := NewInfixExpression(NewIdentifier("i"), OperatorLessThan, NewIntegerLiteral(10))
	update := NewIdentifier("i").IncrPostfix()
	body := NewCodeBlock([]Statement{
		NewAssignmentStatement("sum", 0, NewInfixExpression(NewIdentifier("sum"), OperatorAdd, NewIdentifier("i"))),
	})

	forStmt := NewForStatement(initor, cond, update, body)

	checkInterfaceCodeElement(forStmt)
	checkInterfaceStatement(forStmt)

	expected := strings.Join([]string{
		"for (i = 0; i < 10; i++) {",
		"    sum = sum + i;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle1, expected, forStmt)
}

func TestForStatementWithoutVariableDeclarationStyle2(t *testing.T) {
	initor := NewAssignmentExpression("i", 0, NewIntegerLiteral(0))
	cond := NewInfixExpression(NewIdentifier("i"), OperatorLessThan, NewIntegerLiteral(10))
	update := NewIdentifier("i").IncrPostfix()
	body := NewCodeBlock([]Statement{
		NewAssignmentStatement("sum", 0, NewInfixExpression(NewIdentifier("sum"), OperatorAdd, NewIdentifier("i"))),
	})

	forStmt := NewForStatement(initor, cond, update, body)

	checkInterfaceCodeElement(forStmt)
	checkInterfaceStatement(forStmt)

	expected := strings.Join([]string{
		"for(i = 0; i < 10; i++)",
		"{",
		"    sum = sum + i;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle2, expected, forStmt)
}

func TestSwitchStatementStyle1(t *testing.T) {
	cond := NewIdentifier("a")
	case1Block := NewCodeBlock([]Statement{
		NewAssignmentStatement("r", 0, NewIntegerLiteral(1)),
		NewBreakStatement(),
	})
	case2Block := NewCodeBlock([]Statement{
		NewAssignmentStatement("r", 0, NewIntegerLiteral(2)),
		NewBreakStatement(),
	})
	defaultBlock := NewCodeBlock([]Statement{
		NewAssignmentStatement("r", 0, NewIntegerLiteral(0)),
		NewBreakStatement(),
	})

	switchStmt := NewSwitchStatement(cond, []*CaseBranch{
		NewCaseBranch(NewIntegerLiteral(1), case1Block),
		NewCaseBranch(NewIntegerLiteral(2), case2Block),
	}, defaultBlock)

	checkInterfaceCodeElement(switchStmt)
	checkInterfaceStatement(switchStmt)

	expected := strings.Join([]string{
		"switch (a) {",
		"case 1:",
		"    r = 1;",
		"    break;",
		"case 2:",
		"    r = 2;",
		"    break;",
		"default:",
		"    r = 0;",
		"    break;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle1, expected, switchStmt)
}

func TestSwitchStatementStyle2(t *testing.T) {
	cond := NewIdentifier("a")
	case1Block := NewCodeBlock([]Statement{
		NewAssignmentStatement("r", 0, NewIntegerLiteral(1)),
		NewBreakStatement(),
	})
	case2Block := NewCodeBlock([]Statement{
		NewAssignmentStatement("r", 0, NewIntegerLiteral(2)),
		NewBreakStatement(),
	})
	defaultBlock := NewCodeBlock([]Statement{
		NewAssignmentStatement("r", 0, NewIntegerLiteral(0)),
		NewBreakStatement(),
	})

	switchStmt := NewSwitchStatement(cond, []*CaseBranch{
		NewCaseBranch(NewIntegerLiteral(1), case1Block),
		NewCaseBranch(NewIntegerLiteral(2), case2Block),
	}, defaultBlock)

	checkInterfaceCodeElement(switchStmt)
	checkInterfaceStatement(switchStmt)

	expected := strings.Join([]string{
		"switch(a)",
		"{",
		"case 1:",
		"    r = 1;",
		"    break;",
		"case 2:",
		"    r = 2;",
		"    break;",
		"default:",
		"    r = 0;",
		"    break;",
		"}",
		"",
	}, "\n")
	checkOutputOnStyle(t, testStyle2, expected, switchStmt)
}
