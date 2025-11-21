package ast

import (
	"github.com/flily/magi-c/context"
)

type ReturnStatement struct {
	NonTerminalNode

	Return *TerminalToken
	Value  *IntegerLiteral
}

func (r *ReturnStatement) statementNode() {}

func NewReturnStatement(keyword *TerminalToken) *ReturnStatement {
	s := &ReturnStatement{
		Return: keyword,
	}
	s.Init(s)

	return s
}

func (r *ReturnStatement) Context() *context.Context {
	return context.JoinObjects(r.Return, r.Value)
}
