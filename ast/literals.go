package ast

import "github.com/flily/magi-c/context"

type StringLiteral struct {
	TerminalNode
	Value string
}

func NewStringLiteral(ctx *context.Context, value string) *StringLiteral {
	l := &StringLiteral{
		TerminalNode: NewTerminalNode(ctx),
		Value:        value,
	}

	return l
}
