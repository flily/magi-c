package ast

import (
	"github.com/flily/magi-c/context"
)

type FunctionDeclaration struct {
	NonTerminalNode
	Keyword           *TerminalToken
	Name              *Identifier
	LParenArgs        *TerminalToken
	Arguments         *ArgumentList
	RParenArgs        *TerminalToken
	LParenReturnTypes *TerminalToken
	ReturnTypes       *TypeList
	RParenReturnTypes *TerminalToken
	LBracket          *TerminalToken
	Statements        []Statement
	RBracket          *TerminalToken
}

func (f *FunctionDeclaration) declarationNode() {}

func (f *FunctionDeclaration) EqualTo(other Comparable) bool {
	o, ok := CheckNodeEqual(f, other)
	if !ok {
		return false
	}

	if !f.Keyword.EqualTo(o.Keyword) {
		return false
	}

	if !f.Name.EqualTo(o.Name) {
		return false
	}

	if !f.LParenArgs.EqualTo(o.LParenArgs) {
		return false
	}

	if !CheckNilPointerEqual(f.Arguments, o.Arguments) {
		return false
	}

	if !f.RParenArgs.EqualTo(o.RParenArgs) {
		return false
	}

	if !f.LParenArgs.EqualTo(o.LParenArgs) {
		return false
	}

	if !CheckNilPointerEqual(f.ReturnTypes, o.ReturnTypes) {
		return false
	}

	if !f.RParenReturnTypes.EqualTo(o.RParenReturnTypes) {
		return false
	}

	if !f.LBracket.EqualTo(o.LBracket) {
		return false
	}

	if !CheckArrayEqual(f.Statements, o.Statements) {
		return false
	}

	if !f.RBracket.EqualTo(o.RBracket) {
		return false
	}

	return true
}

func (f *FunctionDeclaration) Context() *context.Context {
	ctx := context.JoinObjects(
		f.Keyword,
		f.Name,
		f.LParenArgs,
		f.Arguments,
		f.RParenArgs,
		f.LParenReturnTypes,
		f.ReturnTypes,
		f.RParenReturnTypes,
		f.LBracket,
		f.RBracket,
	)

	return ctx
}
