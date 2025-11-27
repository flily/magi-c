package ast

import (
	"github.com/flily/magi-c/context"
)

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

func (p *ArgumentDeclaration) EqualTo(other Comparable) bool {
	o, ok := CheckNodeEqual(p, other)
	if !ok {
		return false
	}

	if !p.Name.EqualTo(o.Name) {
		return false
	}

	if !p.Type.EqualTo(o.Type) {
		return false
	}

	if !CheckNilPointerEqual(p.Comma, o.Comma) {
		return false
	}

	return true
}

func (p *ArgumentDeclaration) Context() *context.Context {
	return context.JoinObjects(p.Name, p.Type, p.Comma)
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

func (l *ArgumentList) EqualTo(other Comparable) bool {
	o, ok := CheckNodeEqual(l, other)
	if !ok {
		return false
	}

	return CheckArrayEqual(l.Arguments, o.Arguments)
}

func (l *ArgumentList) Context() *context.Context {
	if len(l.Arguments) == 0 {
		return nil
	}

	ctxList := make([]context.ContextProvider, 0, len(l.Arguments))
	for _, n := range l.Arguments {
		ctxList = append(ctxList, n)
	}

	return context.JoinObjects(ctxList...)
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

func (t *SimpleType) EqualTo(other Comparable) bool {
	o, ok := CheckNodeEqual(t, other)
	if !ok {
		return false
	}

	if !CheckArrayEqual(t.PointerAsterisk, o.PointerAsterisk) {
		return false
	}

	if !t.Identifier.EqualTo(o.Identifier) {
		return false
	}

	return true
}

func (t *SimpleType) Context() *context.Context {
	ctxList := make([]context.ContextProvider, 0, len(t.PointerAsterisk)+1)
	for _, asterisk := range t.PointerAsterisk {
		ctxList = append(ctxList, asterisk)
	}
	ctxList = append(ctxList, t.Identifier)

	return context.JoinObjects(ctxList...)
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

func (t *FunctionType) EqualTo(other Comparable) bool {
	o, ok := CheckNodeEqual(t, other)
	if !ok {
		return false
	}

	if !t.Keyword.EqualTo(o.Keyword) {
		return false
	}

	if !t.ArgumentLParen.EqualTo(o.ArgumentLParen) {
		return false
	}

	if !CheckNilPointerEqual(t.ArgumentList, o.ArgumentList) {
		return false
	}

	if !t.ArgumentRParen.EqualTo(o.ArgumentRParen) {
		return false
	}

	if !CheckNilPointerEqual(t.ReturnTypes, o.ReturnTypes) {
		return false
	}

	if !t.ReturnLParen.EqualTo(o.ReturnLParen) {
		return false
	}

	if !t.ReturnRParen.EqualTo(o.ReturnRParen) {
		return false
	}

	return true
}

func (t *FunctionType) Context() *context.Context {
	ctx := context.JoinObjects(
		t.Keyword,
		t.ArgumentLParen,
		t.ArgumentList,
		t.ArgumentRParen,
		t.ReturnLParen,
		t.ReturnTypes,
		t.ReturnRParen,
	)

	return ctx
}

type TypeListItems struct {
	NonTerminalNode
	Type  Type
	Comma *TerminalToken
}

func NewTypeListItems(t Type) *TypeListItems {
	l := &TypeListItems{
		Type: t,
	}
	l.Init(l)

	return l
}

func (l *TypeListItems) EqualTo(other Comparable) bool {
	o, ok := CheckNodeEqual(l, other)
	if !ok {
		return false
	}

	if !l.Type.EqualTo(o.Type) {
		return false
	}

	if !CheckNilPointerEqual(l.Comma, o.Comma) {
		return false
	}

	return true
}

func (l *TypeListItems) Context() *context.Context {
	return context.JoinObjects(l.Type, l.Comma)
}

type TypeList struct {
	Types []*TypeListItems
}

func NewTypeList() *TypeList {
	l := &TypeList{
		Types: make([]*TypeListItems, 0, 2),
	}

	return l
}

func (l *TypeList) EqualTo(other Comparable) bool {
	o, ok := CheckNodeEqual(l, other)
	if !ok {
		return false
	}

	return CheckArrayEqual(l.Types, o.Types)
}

func (l *TypeList) Context() *context.Context {
	if len(l.Types) == 0 {
		return nil
	}

	ctxList := make([]context.ContextProvider, 0, len(l.Types))
	for _, item := range l.Types {
		ctxList = append(ctxList, item)
	}

	return context.JoinObjects(ctxList...)
}
