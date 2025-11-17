package ast

import "github.com/flily/magi-c/context"

type FunctionDeclaration struct {
	NonTerminalNode

	Keyword           *TerminalToken
	Name              *Identifier
	LParenArgs        *TerminalToken
	Arguments         *ArgumentList
	RParenArgs        *TerminalToken
	LParenReturnTypes *TerminalToken
	RParenReturnTypes *TerminalToken
	LBracket          *TerminalToken
	Statements        []Statement
	RBracket          *TerminalToken
}

func (f *FunctionDeclaration) declarationNode() {}

func (f *FunctionDeclaration) Context() *context.Context {
	ctx := context.Join(
		f.Keyword.Context(),
		f.Name.Context(),
		f.LParenArgs.Context(),
		f.RParenArgs.Context(),
		f.LParenReturnTypes.Context(),
		f.RParenReturnTypes.Context(),
		f.LBracket.Context(),
	)

	if len(f.Statements) == 0 {
		ctx = context.Join(ctx, f.RBracket.Context())
		return ctx
	}

	ctxList := make([]*context.Context, 0, len(f.Statements)+1)
	ctxList = append(ctxList, ctx)
	for _, stmt := range f.Statements {
		ctxList = append(ctxList, stmt.Context())
	}
	ctxList = append(ctxList, f.RBracket.Context())

	return context.Join(ctxList...)
}
