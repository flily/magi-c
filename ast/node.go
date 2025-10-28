package ast

import (
	"github.com/flily/magi-c/context"
)

type Node interface {
	Terminal() bool
	Context() *context.Context
	Type() NodeType
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

type TerminalNode struct {
	ContextContainer
}

func NewTerminalNode(ctx *context.Context) TerminalNode {
	n := TerminalNode{
		ContextContainer: ContextContainer{
			context: ctx,
		},
	}

	return n
}

func (n *TerminalNode) Terminal() bool {
	return true
}

type NonTerminalNode struct {
	ContextProvider
}

func (n *NonTerminalNode) Init(provider ContextProvider) {
	n.ContextProvider = provider
}

func (n *NonTerminalNode) Terminal() bool {
	return false
}

func (c *NonTerminalNode) HighlightText(message string, args ...any) string {
	ctx := c.Context()
	return ctx.HighlightText(message, args...)
}
