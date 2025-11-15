package ast

type ReturnStatement struct {
	NonTerminalNode

	Return *TerminalToken
	Value  *IntegerLiteral
}

func (r *ReturnStatement) Type() NodeType {
	return NodeStatement
}

func (r *ReturnStatement) statementNode() {}
