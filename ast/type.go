package ast

import "github.com/flily/magi-c/context"

type Type interface {
	Node
	typeNode()
}

type ArgumentDeclaration struct {
	NonTerminalNode
	Name  *Identifier
	Type  Type
	Comma *TerminalToken
}

func NewArgumentDeclaration() *ArgumentDeclaration {
	a := &ArgumentDeclaration{}
	a.Init(a)

	return a
}

func (p *ArgumentDeclaration) Context() *context.Context {
	// comma not included
	return context.Join(p.Name.Context(), p.Type.Context())
}

type ArgumentList struct {
	NonTerminalNode
	Arguments []*ArgumentDeclaration
}

func NewArgumentList() *ArgumentList {
	l := &ArgumentList{}
	l.Init(l)

	return l
}

func (l *ArgumentList) Context() *context.Context {
	if len(l.Arguments) == 0 {
		return nil
	}

	ctxList := make([]*context.Context, 0, len(l.Arguments))
	for _, n := range l.Arguments {
		ctxList = append(ctxList, n.Context())
	}

	return context.Join(ctxList...)
}

type SimpleType struct {
	NonTerminalNode
	PointerAsterisk []*TerminalToken
	Identifier      *Identifier
}

func NewSimpleType() *SimpleType {
	t := &SimpleType{
		PointerAsterisk: make([]*TerminalToken, 0, 2),
	}
	t.Init(t)

	return t
}

func (t *SimpleType) typeNode() {}

func (t *SimpleType) Context() *context.Context {
	ctxList := make([]*context.Context, 0, len(t.PointerAsterisk)+1)
	for _, asterisk := range t.PointerAsterisk {
		ctxList = append(ctxList, asterisk.Context())
	}
	ctxList = append(ctxList, t.Identifier.Context())

	return context.Join(ctxList...)
}

func (t *SimpleType) AddPointerAsterisk(asterisk *TerminalToken) {
	t.PointerAsterisk = append(t.PointerAsterisk, asterisk)
}

type FunctionType struct {
	NonTerminalNode
	Keyword        *TerminalToken
	ArgumentLParen *TerminalToken
	ArgumentList   *ArgumentList
	ArgumentRParen *TerminalToken
	ReturnLParen   *TerminalToken
	ReturnTypes    *ArgumentList
	ReturnRParen   *TerminalToken
}

func NewFunctionType() *FunctionType {
	t := &FunctionType{}
	t.Init(t)

	return t
}

func (t *FunctionType) typeNode() {}

func (t *FunctionType) Context() *context.Context {
	ctx := context.Join(
		t.Keyword.Context(),
		t.ArgumentLParen.Context(),
		t.ArgumentList.Context(),
		t.ArgumentRParen.Context(),
		t.ReturnLParen.Context(),
		t.ReturnTypes.Context(),
		t.ReturnRParen.Context(),
	)

	return ctx
}
