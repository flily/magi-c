package ast

import (
	"github.com/flily/magi-c/context"
)

type Node interface {
	Context() *context.Context
	HighlightText(message string, args ...any) string
}

type ContextContainer struct {
	context *context.Context
}

func (n *ContextContainer) Context() *context.Context {
	return n.context
}

func (n *ContextContainer) HighlightText(message string, args ...any) string {
	ctx := n.Context()
	return ctx.HighlightText(message, args...)
}

type TerminalNode interface {
	Node
	Type() TokenType
}

type TerminalNodeBase struct {
	ContextContainer
}

func NewTerminalNodeBase(ctx *context.Context) TerminalNodeBase {
	n := TerminalNodeBase{
		ContextContainer: ContextContainer{
			context: ctx,
		},
	}

	return n
}

type NonTerminalNode struct {
	provider context.ContextProvider
}

func (n *NonTerminalNode) Init(provider context.ContextProvider) {
	n.provider = provider
}

func (c *NonTerminalNode) HighlightText(message string, args ...any) string {
	ctx := c.provider.Context()
	return ctx.HighlightText(message, args...)
}
