package ast

import (
	"github.com/flily/magi-c/context"
)

type Node interface {
	Terminal() bool
	Context() *context.Context
	HighlightText(message string, args ...any) string
}

type ContextProvider interface {
	Context() *context.Context
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

func (n *TerminalNodeBase) Terminal() bool {
	return true
}

type NonTerminalNode struct {
	provider ContextProvider
}

func (n *NonTerminalNode) Init(provider ContextProvider) {
	n.provider = provider
}

func (n *NonTerminalNode) Terminal() bool {
	return false
}

func (c *NonTerminalNode) Context() *context.Context {
	return c.provider.Context()
}

func (c *NonTerminalNode) HighlightText(message string, args ...any) string {
	ctx := c.Context()
	return ctx.HighlightText(message, args...)
}
