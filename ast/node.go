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

type nodeContext struct {
	context *context.Context
}

func (n *nodeContext) Context() *context.Context {
	return n.context
}

func (n *nodeContext) HighlightText(message string, args ...any) string {
	return n.context.HighlightText(message, args...)
}

type TerminalNode struct {
	nodeContext
}

func NewTerminalNode(ctx *context.Context) TerminalNode {
	n := TerminalNode{
		nodeContext: nodeContext{context: ctx},
	}

	return n
}

func (n *TerminalNode) Terminal() bool {
	return true
}

type NonTerminalNode struct {
	nodeContext
}

func NewNonTerminalNode(ctx *context.Context) NonTerminalNode {
	n := NonTerminalNode{
		nodeContext: nodeContext{context: ctx},
	}

	return n
}

func (n *NonTerminalNode) Terminal() bool {
	return false
}
