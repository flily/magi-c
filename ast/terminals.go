package ast

import (
	"github.com/flily/magi-c/context"
)

type StringLiteral struct {
	TerminalNodeBase
	Value string
}

func NewStringLiteral(ctx *context.Context, value string) *StringLiteral {
	l := &StringLiteral{
		TerminalNodeBase: NewTerminalNodeBase(ctx),
		Value:            value,
	}

	return l
}

func (l *StringLiteral) Type() TokenType {
	return String
}

type IntegerLiteral struct {
	TerminalNodeBase
	Value uint64
}

func NewIntegerLiteral(ctx *context.Context, value uint64) *IntegerLiteral {
	l := &IntegerLiteral{
		TerminalNodeBase: NewTerminalNodeBase(ctx),
		Value:            value,
	}

	return l
}

func (l *IntegerLiteral) expressionNode() {}

func (l *IntegerLiteral) Type() TokenType {
	return Integer
}

type FloatLiteral struct {
	TerminalNodeBase
	Value float64
}

func NewFloatLiteral(ctx *context.Context, value float64) *FloatLiteral {
	l := &FloatLiteral{
		TerminalNodeBase: NewTerminalNodeBase(ctx),
		Value:            value,
	}

	return l
}

func (l *FloatLiteral) Type() TokenType {
	return Float
}

type Identifier struct {
	TerminalNodeBase
	Name string
}

func NewIdentifier(ctx *context.Context, name string) *Identifier {
	id := &Identifier{
		TerminalNodeBase: NewTerminalNodeBase(ctx),
		Name:             name,
	}

	return id
}

func (i *Identifier) Type() TokenType {
	return IdentifierName
}

type TerminalToken struct {
	TerminalNodeBase
	Token TokenType
}

func NewTerminalToken(ctx *context.Context, token TokenType) *TerminalToken {
	t := &TerminalToken{
		TerminalNodeBase: NewTerminalNodeBase(ctx),
		Token:            token,
	}

	return t
}

func (t *TerminalToken) Type() TokenType {
	return t.Token
}
