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

func NewFunctionDeclaration(keyword *TerminalToken) *FunctionDeclaration {
	f := &FunctionDeclaration{
		Keyword: keyword,
	}
	f.Init(f)

	return f
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
	funcDecl.Init(funcDecl)

	return funcDecl
}

func (f *FunctionDeclaration) declarationNode() {}

func (f *FunctionDeclaration) EqualTo(archor context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(f, other)
	if err != nil {
		return err
	}

	if err := f.Name.EqualTo(f, o.Name); err != nil {
		return err
	}

	if err := CheckNilPointerEqual(f, f.Arguments, o.Arguments); err != nil {
		return err
	}

	if err := CheckNilPointerEqual(f, f.ReturnTypes, o.ReturnTypes); err != nil {
		return err
	}

	if err := CheckArrayEqual("STATEMENT LIST", f, f.Statements, o.Statements); err != nil {
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
