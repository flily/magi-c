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
