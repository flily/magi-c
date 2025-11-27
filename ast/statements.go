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

func (r *ReturnStatement) EqualTo(other Comparable) bool {
	o, ok := CheckNodeEqual(r, other)
	if !ok {
		return false
	}

	return CheckNilPointerEqual(r.Value, o.Value)
}

func (r *ReturnStatement) Context() *context.Context {
	return context.JoinObjects(r.Return, r.Value)
}
