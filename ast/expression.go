package ast

import (
	"github.com/flily/magi-c/context"
)

type Expression interface {
	Node
	expressionNode()
}

type ExpressionListItem struct {
	NonTerminalNode
	Expression Expression
	Comma      *TerminalToken
}

func NewExpressionListItem(e Expression, comma *TerminalToken) *ExpressionListItem {
	item := &ExpressionListItem{
		Expression: e,
		Comma:      comma,
	}
	item.Init(item)

	return item
}

func ASTBuildExpressionListItem(e Expression, comma *TerminalToken) *ExpressionListItem {
	return NewExpressionListItem(e, comma)
}

func (i *ExpressionListItem) EqualTo(other Comparable) error {
	o, err := CheckNodeEqual(i, other)
	if err != nil {
		return err
	}

	if err := i.Expression.EqualTo(o.Expression); err != nil {
		return err
	}

	if err := CheckNilPointerEqual(i, i.Comma, o.Comma); err != nil {
		return err
	}

	return nil
}

func (i *ExpressionListItem) Context() *context.Context {
	return context.JoinObjects(i.Expression, i.Comma)
}

type ExpressionList struct {
	Expressions []*ExpressionListItem
}

func NewExpressionList() *ExpressionList {
	l := &ExpressionList{
		Expressions: make([]*ExpressionListItem, 0, 2),
	}

	return l
}

func ASTBuildExpressionList(items ...*ExpressionListItem) *ExpressionList {
	l := &ExpressionList{
		Expressions: items,
	}

	return l
}

func (l *ExpressionList) EqualTo(other Comparable) error {
	o, err := CheckNodeEqual(l, other)
	if err != nil {
		return err
	}

	return CheckArrayEqual(l.Expressions, o.Expressions)
}

func (l *ExpressionList) Context() *context.Context {
	if l == nil || len(l.Expressions) == 0 {
		return nil
	}

	ctxList := make([]context.ContextProvider, 0, len(l.Expressions))
	for _, item := range l.Expressions {
		ctxList = append(ctxList, item)
	}

	return context.JoinObjects(ctxList...)
}

type InfixExpression struct {
	NonTerminalNode
	LeftOperand  Expression
	Operator     *TerminalToken
	RightOperand Expression
}

func NewInfixExpression(left Expression, operator *TerminalToken, right Expression) *InfixExpression {
	expr := &InfixExpression{
		LeftOperand:  left,
		Operator:     operator,
		RightOperand: right,
	}
	expr.Init(expr)

	return expr
}

func (e *InfixExpression) expressionNode() {}

func (e *InfixExpression) EqualTo(other Comparable) error {
	o, err := CheckNodeEqual(e, other)
	if err != nil {
		return err
	}

	if err := e.LeftOperand.EqualTo(o.LeftOperand); err != nil {
		return err
	}

	if e.Operator.Token != o.Operator.Token {
		return NewError(e.Operator.Context(), "expected operator %q, got %q", o.Operator.Token, e.Operator.Token)
	}

	if err := e.RightOperand.EqualTo(o.RightOperand); err != nil {
		return err
	}

	return nil
}

func (e *InfixExpression) Context() *context.Context {
	return context.JoinObjects(e.LeftOperand, e.Operator, e.RightOperand)
}
