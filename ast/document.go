package ast

import (
	"github.com/flily/magi-c/context"
)

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Document struct {
	Statements []Statement
}

func NewDocument(statements []Statement) *Document {
	d := &Document{
		Statements: statements,
	}

	return d
}

func (d *Document) statementNode() {}

func (d *Document) Terminal() bool {
	return false
}

func (d *Document) Type() NodeType {
	return NodeDocument
}

func (d *Document) Context() *context.Context {
	if len(d.Statements) == 0 {
		return nil
	}

	ctxList := make([]*context.Context, 0, len(d.Statements))
	for _, n := range d.Statements {
		ctxList = append(ctxList, n.Context())
	}

	return context.Join(ctxList...)
}

func (d *Document) HighlightText(message string, args ...any) string {
	ctx := d.Context()
	return ctx.HighlightText(message, args...)
}
