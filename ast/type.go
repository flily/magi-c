package ast

import (
	"github.com/flily/magi-c/context"
)

type Type interface {
	Node
	typeNode()
}

type SimpleType struct {
	NonTerminalNode
	PointerAsterisk []*TerminalToken
	Identifier      *Identifier
}

func NewSimpleType(asterisks []*TerminalToken, identifier *Identifier) *SimpleType {
	t := &SimpleType{
		PointerAsterisk: asterisks,
		Identifier:      identifier,
	}
	t.Init(t)

	return t
}

func ASTBuildSimpleType(name string) *SimpleType {
	t := NewSimpleType(nil, nil)

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

type ArgumentDeclaration struct {
	NonTerminalNode
	Name  *Identifier
	Type  Type
	Comma *TerminalToken
}

func NewArgumentDeclaration(name *Identifier, t Type, comma *TerminalToken) *ArgumentDeclaration {
	a := &ArgumentDeclaration{
		Name:  name,
		Type:  t,
		Comma: comma,
	}
	a.Init(a)

	return a
}

func ASTBuildArgumentWithComma(name string, t string) *ArgumentDeclaration {
	a := NewArgumentDeclaration(ASTBuildIdentifier(name), ASTBuildSimpleType(t), ASTBuildSymbol(Comma))
	return a
}

func ASTBuildArgumentWithoutComma(name string, t string) *ArgumentDeclaration {
	a := NewArgumentDeclaration(ASTBuildIdentifier(name), ASTBuildSimpleType(t), nil)
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

func (l *ArgumentList) Add(name *Identifier, t Type, comma *TerminalToken) {
	arg := &ArgumentDeclaration{
		Name:  name,
		Type:  t,
		Comma: comma,
	}
	l.Arguments = append(l.Arguments, arg)
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

func ASTBuildTypeListItemWithComma(t string) *TypeListItem {
	return NewTypeListItem(ASTBuildSimpleType(t), ASTBuildSymbol(Comma))
}

func ASTBuildTypeListItemWithoutComma(t string) *TypeListItem {
	return NewTypeListItem(ASTBuildSimpleType(t), nil)
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

func (l *TypeList) Add(t Type, comma *TerminalToken) {
	item := NewTypeListItem(t, comma)
	l.Types = append(l.Types, item)
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
