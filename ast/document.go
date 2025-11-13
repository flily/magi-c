package ast

import (
	"github.com/flily/magi-c/context"
)

type Declaration interface {
	Node
	declarationNode()
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Document struct {
	Filename     string
	Declarations []Declaration
}

func NewDocument(declarations []Declaration) *Document {
	d := &Document{
		Declarations: declarations,
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
	if len(d.Declarations) == 0 {
		return nil
	}

	ctxList := make([]*context.Context, 0, len(d.Declarations))
	for _, n := range d.Declarations {
		ctxList = append(ctxList, n.Context())
	}

	return context.Join(ctxList...)
}

func (d *Document) HighlightText(message string, args ...any) string {
	ctx := d.Context()
	return ctx.HighlightText(message, args...)
}
