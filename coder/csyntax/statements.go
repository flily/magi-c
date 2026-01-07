package csyntax

import (
	"strings"
)

type BlockContext int

const (
	BlockContextFunction BlockContext = iota
	BlockContextIf
	BlockContextElse
	BlockContextDo
	BlockContextWhile
	BlockContextFor
	BlockContextSwitch
	BlockContextCase
	BlockContextDefault
)

type CodeBlock struct {
	Statements []Statement
}

func NewCodeBlock(statements []Statement) *CodeBlock {
	b := &CodeBlock{
		Statements: statements,
	}

	return b
}

func (b *CodeBlock) codeElement() {}

func (b *CodeBlock) statementNode() {}

func (b *CodeBlock) Write(out *StyleWriter, level int) error {
	if err := out.WriteIndent(level); err != nil {
		return err
	}

	if err := out.WriteLine(LeftBrace); err != nil {
		return err
	}

	for _, stmt := range b.Statements {
		if err := stmt.Write(out, level+1); err != nil {
			return err
		}
	}

	if err := out.WriteIndent(level); err != nil {
		return err
	}

	return out.WriteLine(RightBrace)
}

type AssignmentStatement struct {
	LeftIdentifier   string
	LeftPointerLevel int
	RightExpression  Expression
}

func NewAssignmentStatement(leftIdentifier string, leftPointerLevel int, rightExpression Expression) *AssignmentStatement {
	s := &AssignmentStatement{
		LeftIdentifier:   leftIdentifier,
		LeftPointerLevel: leftPointerLevel,
		RightExpression:  rightExpression,
	}

	return s
}

func (s *AssignmentStatement) codeElement() {}

func (s *AssignmentStatement) statementNode() {}

func (s *AssignmentStatement) Write(out *StyleWriter, level int) error {
	if err := out.WriteIndent(level); err != nil {
		return err
	}

	pointer := strings.Repeat(PointerAsterisk, s.LeftPointerLevel)
	pointerSpace := ""
	if s.LeftPointerLevel > 0 && out.style.PointerSpacingBefore {
		pointerSpace = Space
	}

	if err := out.Write("%s%s%s%s", pointer, pointerSpace, s.LeftIdentifier, out.style.Assign()); err != nil {
		return err
	}

	if err := s.RightExpression.Write(out, level); err != nil {
		return err
	}

	return out.WriteLine(Semicolon)
}

type ReturnStatement struct {
	Expression Expression
}

func NewReturnStatement(expression Expression) *ReturnStatement {
	s := &ReturnStatement{
		Expression: expression,
	}

	return s
}

func (s *ReturnStatement) statementNode() {}

func (s *ReturnStatement) codeElement() {}

func (s *ReturnStatement) Write(out *StyleWriter, level int) error {
	if err := out.WriteIndent(level); err != nil {
		return err
	}

	if s.Expression == nil {
		return out.WriteItems(level, KeywordReturn, PunctuatorSemicolon, out.EOL) // return;
	}

	return out.WriteItems(level, KeywordReturn, DelimiterSpace, s.Expression, PunctuatorSemicolon, out.EOL) // return expr;
}
