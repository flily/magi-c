package ast

import (
	"github.com/flily/magi-c/context"
)

type Document struct {
	Nodes []Node
}

func NewDocument(nodes []Node) *Document {
	d := &Document{
		Nodes: nodes,
	}

	return d
}

func (d *Document) Terminal() bool {
	return false
}

func (d *Document) Type() NodeType {
	return NodeDocument
}

func (d *Document) Context() *context.Context {
	if len(d.Nodes) == 0 {
		return nil
	}

	ctxList := make([]*context.Context, 0, len(d.Nodes))
	for _, n := range d.Nodes {
		ctxList = append(ctxList, n.Context())
	}

	return context.Join(ctxList...)
}
