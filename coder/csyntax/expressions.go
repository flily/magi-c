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

func (id *Identifier) Write(out *StyleWriter, level Level) error {
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

func (e *InfixExpression) Write(out *StyleWriter, level Level) error {
	return out.Write(level.NextParanthesis(),
		OperatorLeftParen.Select(level.ParanthesisLevel > 0),
		e.Left, out.style.BinaryOperator(e.Operator), e.Right,
		OperatorRightParen.Select(level.ParanthesisLevel > 0))
}

type PostfixExpression struct {
	Operand  Expression
	Operator Punctuator
}

func NewPostfixExpression(operand Expression, operator Punctuator) *PostfixExpression {
	expr := &PostfixExpression{
		Operand:  operand,
		Operator: operator,
	}

	return expr
}

func (e *PostfixExpression) codeElement()    {}
func (e *PostfixExpression) expressionNode() {}

func (e *PostfixExpression) Write(out *StyleWriter, level Level) error {
	return out.Write(level.NextParanthesis(),
		OperatorLeftParen.Select(level.ParanthesisLevel > 0),
		e.Operand, e.Operator,
		OperatorRightParen.Select(level.ParanthesisLevel > 0),
	)
}

type UnaryExpression struct {
	Operator Punctuator
	Operand  Expression
}

func NewUnaryExpression(operator Punctuator, operand Expression) *UnaryExpression {
	expr := &UnaryExpression{
		Operator: operator,
		Operand:  operand,
	}

	return expr
}

func (e *UnaryExpression) codeElement()    {}
func (e *UnaryExpression) expressionNode() {}

func (e *UnaryExpression) Write(out *StyleWriter, level Level) error {
	return out.Write(level.NextParanthesis(),
		OperatorLeftParen.Select(level.ParanthesisLevel > 0),
		e.Operator, e.Operand,
		OperatorRightParen.Select(level.ParanthesisLevel > 0),
	)
}
