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

func ASTBuildArgumentWithComma(name string, t Type) *ArgumentDeclaration {
	a := NewArgumentDeclaration()
	a.Name = ASTBuildIdentifier(name)
	a.Type = t
	a.Comma = ASTBuildSymbol(Comma)
	return a
}

func ASTBuildArgumentWithoutComma(name string, t Type) *ArgumentDeclaration {
	a := NewArgumentDeclaration()
	a.Name = ASTBuildIdentifier(name)
	a.Type = t
	return a
}

func (p *ArgumentDeclaration) EqualTo(_ context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(p, other)
	if err != nil {
		return err
	}

	if err := p.Name.EqualTo(p, o.Name); err != nil {
		return err
	}

	if err := p.Type.EqualTo(p, o.Type); err != nil {
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

func ASTBuildArgumentList(args ...*ArgumentDeclaration) *ArgumentList {
	l := NewArgumentList()
	l.Arguments = args

	return l
}

func (l *ArgumentList) EqualTo(archor context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(l, other)
	if err != nil {
		return err
	}

	return CheckArrayEqual("ARGUMENT LIST", archor, l.Arguments, o.Arguments)
}

func (l *ArgumentList) Context() *context.Context {
	if l == nil || len(l.Arguments) == 0 {
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

func (t *SimpleType) EqualTo(_ context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(t, other)
	if err != nil {
		return err
	}

	if err := CheckArrayEqual("POINTER ASTERISK", t, t.PointerAsterisk, o.PointerAsterisk); err != nil {
		return err
	}

	if err := t.Identifier.EqualTo(t, o.Identifier); err != nil {
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

func (t *FunctionType) EqualTo(_ context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(t, other)
	if err != nil {
		return err
	}

	if err := t.Keyword.EqualTo(t, o.Keyword); err != nil {
		return err
	}

	if err := t.ArgumentLParen.EqualTo(t, o.ArgumentLParen); err != nil {
		return err
	}

	if err := CheckNilPointerEqual(t, t.ArgumentList, o.ArgumentList); err != nil {
		return err
	}

	if err := t.ArgumentRParen.EqualTo(t, o.ArgumentRParen); err != nil {
		return err
	}

	if err := CheckNilPointerEqual(t, t.ReturnTypes, o.ReturnTypes); err != nil {
		return err
	}

	if err := t.ReturnLParen.EqualTo(t, o.ReturnLParen); err != nil {
		return err
	}

	if err := t.ReturnRParen.EqualTo(t, o.ReturnRParen); err != nil {
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

func ASTBuildTypeListItemWithComma(t Type) *TypeListItem {
	return NewTypeListItem(t, ASTBuildSymbol(Comma))
}

func ASTBuildTypeListItemWithoutComma(t Type) *TypeListItem {
	return NewTypeListItem(t, nil)
}

func (l *TypeListItem) EqualTo(_ context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(l, other)
	if err != nil {
		return err
	}

	if err := l.Type.EqualTo(l, o.Type); err != nil {
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

func (l *TypeList) EqualTo(archor context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(l, other)
	if err != nil {
		return err
	}

	return CheckArrayEqual("TYPE LIST", archor, l.Types, o.Types)
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
