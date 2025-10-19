package ast

import (
	"github.com/flily/magi-c/context"
)

type Token struct {
	t   NodeType
	ctx *context.Context
}

func NewToken(ttype NodeType, ctx *context.Context) *Token {
	t := &Token{
		t:   ttype,
		ctx: ctx,
	}

	return t
}

func (t *Token) Type() NodeType {
	return t.t
}

func (t *Token) Context() *context.Context {
	return t.ctx
}

func (t *Token) HighlightText(s string) string {
	return t.ctx.HighlightText(s)
}
