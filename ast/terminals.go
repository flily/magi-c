package ast

import (
	"github.com/flily/magi-c/context"
)

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

func (l *StringLiteral) Type() TokenType {
	return String
}

type IntegerLiteral struct {
	TerminalNode
	Value uint64
}

func NewIntegerLiteral(ctx *context.Context, value uint64) *IntegerLiteral {
	l := &IntegerLiteral{
		TerminalNode: NewTerminalNode(ctx),
		Value:        value,
	}

	return l
}

func (l *IntegerLiteral) Type() TokenType {
	return Integer
}

type FloatLiteral struct {
	TerminalNode
	Value float64
}

func NewFloatLiteral(ctx *context.Context, value float64) *FloatLiteral {
	l := &FloatLiteral{
		TerminalNode: NewTerminalNode(ctx),
		Value:        value,
	}

	return l
}

func (l *FloatLiteral) Type() TokenType {
	return Float
}

type Identifier struct {
	TerminalNode
	Name string
}

func NewIdentifier(ctx *context.Context, name string) *Identifier {
	id := &Identifier{
		TerminalNode: NewTerminalNode(ctx),
		Name:         name,
	}

	return id
}

func (i *Identifier) Type() TokenType {
	return IdentifierName
}

type TerminalToken struct {
	TerminalNode
	Token TokenType
}

func NewTerminalToken(ctx *context.Context, token TokenType) *TerminalToken {
	t := &TerminalToken{
		TerminalNode: NewTerminalNode(ctx),
		Token:        token,
	}

	return t
}

func (t *TerminalToken) Type() TokenType {
	return t.Token
}
