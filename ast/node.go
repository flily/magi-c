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

func NewNonTerminalNode(provider ContextProvider) NonTerminalNode {
	n := NonTerminalNode{
		ContextProvider: provider,
	}

	return n
}

func (n *NonTerminalNode) Terminal() bool {
	return false
}
