package csyntax

type InfixExpression struct {
	Left     Expression
	Operator CodeElement
	Right    Expression
}
