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

type IntegerLiteral struct {
	TerminalNode
	Type  BaseType
	Sign  int
	Value uint64
}

func NewIntegerLiteral(ctx *context.Context, sign int, value uint64) *IntegerLiteral {
	l := &IntegerLiteral{
		TerminalNode: NewTerminalNode(ctx),
		Sign:         sign,
		Value:        value,
	}

	return l
}

type FloatLiteral struct {
	TerminalNode
	Sign  int
	Value float64
}

func NewFloatLiteral(ctx *context.Context, sign int, value float64) *FloatLiteral {
	l := &FloatLiteral{
		TerminalNode: NewTerminalNode(ctx),
		Sign:         sign,
		Value:        value,
	}

	return l
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
