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

func ASTBuildFunction(name string, args []*ArgumentDeclaration, returnTypes *TypeList, statements []Statement) *FunctionDeclaration {
	funcDecl := &FunctionDeclaration{
		Keyword:           NewTerminalToken(nil, Function),
		Name:              ASTBuildIdentifier(name),
		LParenArgs:        NewTerminalToken(nil, LeftParen),
		Arguments:         nil,
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

func (f *FunctionDeclaration) EqualTo(other Comparable) error {
	o, err := CheckNodeEqual(f, other)
	if err != nil {
		return err
	}

	if err := f.Keyword.EqualTo(o.Keyword); err != nil {
		return err
	}

	if err := f.Name.EqualTo(o.Name); err != nil {
		return err
	}

	if err := f.LParenArgs.EqualTo(o.LParenArgs); err != nil {
		return err
	}

	if err := CheckNilPointerEqual(f, f.Arguments, o.Arguments); err != nil {
		return err
	}

	if err := f.RParenArgs.EqualTo(o.RParenArgs); err != nil {
		return err
	}

	if err := f.LParenArgs.EqualTo(o.LParenArgs); err != nil {
		return err
	}

	if err := CheckNilPointerEqual(f, f.ReturnTypes, o.ReturnTypes); err != nil {
		return err
	}

	if err := f.RParenReturnTypes.EqualTo(o.RParenReturnTypes); err != nil {
		return err
	}

	if err := f.LBrace.EqualTo(o.LBrace); err != nil {
		return err
	}

	if err := CheckArrayEqual(f.Statements, o.Statements); err != nil {
		return err
	}

	if err := f.RBrace.EqualTo(o.RBrace); err != nil {
		return err
	}

	return nil
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
		f.LBrace,
		f.RBrace,
	)

	return ctx
}
