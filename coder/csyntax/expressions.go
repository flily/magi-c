package csyntax

type Identifier struct {
	Name string
}

func NewIdentifier(name string) *Identifier {
	id := &Identifier{
		Name: name,
	}

	return id
}

func (id *Identifier) codeElement()    {}
func (id *Identifier) expressionNode() {}

func (id *Identifier) Write(out *StyleWriter, level int) error {
	return out.Write(level, StringElement(id.Name))
}

type InfixExpression struct {
	Left     Expression
	Operator Punctuator
	Right    Expression
}

func NewInfixExpression(left Expression, operator Punctuator, right Expression) *InfixExpression {
	expr := &InfixExpression{
		Left:     left,
		Operator: operator,
		Right:    right,
	}

	return expr
}

func (e *InfixExpression) codeElement()    {}
func (e *InfixExpression) expressionNode() {}

func (e *InfixExpression) Write(out *StyleWriter, level int) error {
	return out.Write(level, OperatorLeftParen, e.Left, out.style.BinaryOperator(e.Operator), e.Right, OperatorRightParen)
}
