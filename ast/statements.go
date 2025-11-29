package ast

import (
	"github.com/flily/magi-c/context"
)

type ReturnStatement struct {
	NonTerminalNode

	Return *TerminalToken
	Value  *ExpressionList
}

func (r *ReturnStatement) statementNode() {}

func NewReturnStatement(keyword *TerminalToken) *ReturnStatement {
	s := &ReturnStatement{
		Return: keyword,
	}
	s.Init(s)

	return s
}

func ASTBuildReturnStatement(value *ExpressionList) *ReturnStatement {
	s := NewReturnStatement(ASTBuildKeyword(Return))
	s.Value = value

	return s
}

func (r *ReturnStatement) EqualTo(other Comparable) error {
	o, err := CheckNodeEqual(r, other)
	if err != nil {
		return err
	}

	return CheckNilPointerEqual(r, r.Value, o.Value)
}

func (r *ReturnStatement) Context() *context.Context {
	return context.JoinObjects(r.Return, r.Value)
}
