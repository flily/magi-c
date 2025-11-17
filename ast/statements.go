package ast

type ReturnStatement struct {
	NonTerminalNode

	Return *TerminalToken
	Value  *IntegerLiteral
}

func (r *ReturnStatement) statementNode() {}
