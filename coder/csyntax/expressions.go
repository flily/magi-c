package csyntax

type Expression interface {
	Node
	ForInitializer
	expressionNode()

	IncrPrefix() Expression
	DecrPrefix() Expression
	IncrPostfix() Expression
	DecrPostfix() Expression
}

type ExpressionBase[T Expression] struct {
	Expression T
}

func (e *ExpressionBase[T]) Init(expr T) T {
	e.Expression = expr
	return expr
}

func (e *ExpressionBase[T]) forInitializerNode() {}

func (e *ExpressionBase[T]) IncrPrefix() Expression {
	return NewUnaryExpression(OperatorIncrement, e.Expression)
}

func (e *ExpressionBase[T]) DecrPrefix() Expression {
	return NewUnaryExpression(OperatorDecrement, e.Expression)
}

func (e *ExpressionBase[T]) IncrPostfix() Expression {
	return NewPostfixExpression(e.Expression, OperatorIncrement)
}

func (e *ExpressionBase[T]) DecrPostfix() Expression {
	return NewPostfixExpression(e.Expression, OperatorDecrement)
}

func (e *ExpressionBase[T]) Add(other Expression) Expression {
	return NewInfixExpression(e.Expression, OperatorAdd, other)
}

func (e *ExpressionBase[T]) Sub(other Expression) Expression {
	return NewInfixExpression(e.Expression, OperatorSubtract, other)
}

func (e *ExpressionBase[T]) Mul(other Expression) Expression {
	return NewInfixExpression(e.Expression, OperatorMultiply, other)
}

func (e *ExpressionBase[T]) Div(other Expression) Expression {
	return NewInfixExpression(e.Expression, OperatorDivide, other)
}

type Identifier struct {
	ExpressionBase[*Identifier]
	Name string
}

func NewIdentifier(name string) *Identifier {
	id := &Identifier{
		Name: name,
	}

	return id.Init(id)
}

func (id *Identifier) codeElement()    {}
func (id *Identifier) expressionNode() {}

func (id *Identifier) Write(out *StyleWriter, level Level) error {
	return out.Write(level, StringElement(id.Name))
}

type InfixExpression struct {
	ExpressionBase[*InfixExpression]
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

	return expr.Init(expr)
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
	ExpressionBase[*PostfixExpression]
	Operand  Expression
	Operator Punctuator
}

func NewPostfixExpression(operand Expression, operator Punctuator) *PostfixExpression {
	expr := &PostfixExpression{
		Operand:  operand,
		Operator: operator,
	}

	return expr.Init(expr)
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
	ExpressionBase[*UnaryExpression]
	Operator Punctuator
	Operand  Expression
}

func NewUnaryExpression(operator Punctuator, operand Expression) *UnaryExpression {
	expr := &UnaryExpression{
		Operator: operator,
		Operand:  operand,
	}

	return expr.Init(expr)
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
