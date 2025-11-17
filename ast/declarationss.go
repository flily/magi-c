package ast

type FunctionDeclaration struct {
	NonTerminalNode

	Keyword           *TerminalToken
	Name              *Identifier
	LParenArgs        *TerminalToken
	RParenArgs        *TerminalToken
	LParenReturnTypes *TerminalToken
	RParenReturnTypes *TerminalToken
	LBracket          *TerminalToken
	Statements        []Statement
	RBracket          *TerminalToken
}

func (f *FunctionDeclaration) declarationNode() {}
