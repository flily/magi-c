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
	LBrace            *TerminalToken
	Statements        []Statement
	RBrace            *TerminalToken
}

func ASTBuildFunction(name string, args *ArgumentList, returnTypes *TypeList, statements []Statement) *FunctionDeclaration {
	funcDecl := &FunctionDeclaration{
		Keyword:           NewTerminalToken(nil, Function),
		Name:              ASTBuildIdentifier(name),
		LParenArgs:        NewTerminalToken(nil, LeftParen),
		Arguments:         args,
		RParenArgs:        NewTerminalToken(nil, RightParen),
		LParenReturnTypes: NewTerminalToken(nil, LeftParen),
		ReturnTypes:       returnTypes,
		RParenReturnTypes: NewTerminalToken(nil, RightParen),
		LBrace:            NewTerminalToken(nil, LeftBrace),
		Statements:        statements,
		RBrace:            NewTerminalToken(nil, RightBrace),
	}

	return funcDecl
}

func (f *FunctionDeclaration) declarationNode() {}

func (f *FunctionDeclaration) EqualTo(archor context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(f, other)
	if err != nil {
		return err
	}

	if err := f.Keyword.EqualTo(f, o.Keyword); err != nil {
		return err
	}

	if err := f.Name.EqualTo(f, o.Name); err != nil {
		return err
	}

	if err := f.LParenArgs.EqualTo(f, o.LParenArgs); err != nil {
		return err
	}

	if err := CheckNilPointerEqual(f, f.Arguments, o.Arguments); err != nil {
		return err
	}

	if err := f.RParenArgs.EqualTo(f, o.RParenArgs); err != nil {
		return err
	}

	if err := f.LParenArgs.EqualTo(f, o.LParenArgs); err != nil {
		return err
	}

	if err := CheckNilPointerEqual(f, f.ReturnTypes, o.ReturnTypes); err != nil {
		return err
	}

	if err := f.RParenReturnTypes.EqualTo(f, o.RParenReturnTypes); err != nil {
		return err
	}

	if err := f.LBrace.EqualTo(f, o.LBrace); err != nil {
		return err
	}

	if err := CheckArrayEqual("STATEMENT LIST", f, f.Statements, o.Statements); err != nil {
		return err
	}

	if err := f.RBrace.EqualTo(f, o.RBrace); err != nil {
		return err
	}

	return nil
}

func (f *FunctionDeclaration) Context() *context.Context {
	ctx1 := context.JoinObjects(
		f.Keyword,
		f.Name,
		f.LParenArgs,
		f.Arguments,
		f.RParenArgs,
		f.LParenReturnTypes,
		f.ReturnTypes,
		f.RParenReturnTypes,
		f.LBrace,
	)

	stmtListCtxs := make([]context.ContextProvider, 0, len(f.Statements))
	for _, stmt := range f.Statements {
		stmtListCtxs = append(stmtListCtxs, stmt)
	}

	ctx2 := context.JoinObjects(stmtListCtxs...)
	ctx3 := context.JoinObjects(f.RBrace)

	return context.Join(ctx1, ctx2, ctx3)

}
