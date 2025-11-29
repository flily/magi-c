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

func (p *ArgumentDeclaration) EqualTo(other Comparable) error {
	o, err := CheckNodeEqual(p, other)
	if err != nil {
		return err
	}

	if err := p.Name.EqualTo(o.Name); err != nil {
		return err
	}

	if err := p.Type.EqualTo(o.Type); err != nil {
		return err
	}

	if err := CheckNilPointerEqual(p, p.Comma, o.Comma); err != nil {
		return err
	}

	return nil
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

func (l *ArgumentList) EqualTo(other Comparable) error {
	o, err := CheckNodeEqual(l, other)
	if err != nil {
		return err
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

func ASTBuildSimpleType(name string) *SimpleType {
	t := NewSimpleType()

	start := 0
	for i, c := range name {
		if c == '*' {
			asterisk := NewTerminalToken(nil, Asterisk)
			t.AddPointerAsterisk(asterisk)
			start = i + 1

		} else {
			break
		}
	}

	t.Identifier = ASTBuildIdentifier(name[start:])
	return t
}

func (t *SimpleType) typeNode() {}

func (t *SimpleType) EqualTo(other Comparable) error {
	o, err := CheckNodeEqual(t, other)
	if err != nil {
		return err
	}

	if err := CheckArrayEqual(t.PointerAsterisk, o.PointerAsterisk); err != nil {
		return err
	}

	if err := t.Identifier.EqualTo(o.Identifier); err != nil {
		return err
	}

	return nil
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

func (t *FunctionType) EqualTo(other Comparable) error {
	o, err := CheckNodeEqual(t, other)
	if err != nil {
		return err
	}

	if err := t.Keyword.EqualTo(o.Keyword); err != nil {
		return err
	}

	if err := t.ArgumentLParen.EqualTo(o.ArgumentLParen); err != nil {
		return err
	}

	if err := CheckNilPointerEqual(t, t.ArgumentList, o.ArgumentList); err != nil {
		return err
	}

	if err := t.ArgumentRParen.EqualTo(o.ArgumentRParen); err != nil {
		return err
	}

	if err := CheckNilPointerEqual(t, t.ReturnTypes, o.ReturnTypes); err != nil {
		return err
	}

	if err := t.ReturnLParen.EqualTo(o.ReturnLParen); err != nil {
		return err
	}

	if err := t.ReturnRParen.EqualTo(o.ReturnRParen); err != nil {
		return err
	}

	return nil
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

type TypeListItem struct {
	NonTerminalNode
	Type  Type
	Comma *TerminalToken
}

func NewTypeListItem(t Type, comma *TerminalToken) *TypeListItem {
	l := &TypeListItem{
		Type:  t,
		Comma: comma,
	}
	l.Init(l)

	return l
}

func ASTBuildTypeListItem(t Type, comma *TerminalToken) *TypeListItem {
	return NewTypeListItem(t, comma)
}

func (l *TypeListItem) EqualTo(other Comparable) error {
	o, err := CheckNodeEqual(l, other)
	if err != nil {
		return err
	}

	if err := l.Type.EqualTo(o.Type); err != nil {
		return err
	}

	if err := CheckNilPointerEqual(l, l.Comma, o.Comma); err != nil {
		return err
	}

	return nil
}

func (l *TypeListItem) Context() *context.Context {
	return context.JoinObjects(l.Type, l.Comma)
}

type TypeList struct {
	NonTerminalNode
	Types []*TypeListItem
}

func NewTypeList(items ...*TypeListItem) *TypeList {
	l := &TypeList{
		Types: items,
	}
	l.Init(l)

	return l
}

func ASTBuildTypeList(items ...*TypeListItem) *TypeList {
	return NewTypeList(items...)
}

func (l *TypeList) EqualTo(other Comparable) error {
	o, err := CheckNodeEqual(l, other)
	if err != nil {
		return err
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
