package ast

type Type interface {
	Node
}

type ArgumentList struct {
	NonTerminalNode
	Arguments []*ArgumentDeclaration
}

type SimpleType struct {
	NonTerminalNode
	PointerAsterisk *TerminalToken
	Identifier      *Identifier
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
