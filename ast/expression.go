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

func (i *ExpressionListItem) EqualTo(other Comparable) bool {
	o, ok := CheckNodeEqual(i, other)
	if !ok {
		return false
	}

	if !i.Expression.EqualTo(o.Expression) {
		return false
	}

	return CheckNilPointerEqual(i.Comma, o.Comma)
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

func (l *ExpressionList) EqualTo(other Comparable) bool {
	o, ok := CheckNodeEqual(l, other)
	if !ok {
		return false
	}

	return CheckArrayEqual(l.Expressions, o.Expressions)
}

func (l *ExpressionList) Context() *context.Context {
	if len(l.Expressions) == 0 {
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

func (e *InfixExpression) EqualTo(other Comparable) bool {
	o, ok := CheckNodeEqual(e, other)
	if !ok {
		return false
	}

	if !e.LeftOperand.EqualTo(o.LeftOperand) {
		return false
	}

	if e.Operator.Token != o.Operator.Token {
		return false
	}

	if !e.RightOperand.EqualTo(o.RightOperand) {
		return false
	}

	return true
}

func (e *InfixExpression) Context() *context.Context {
	return context.JoinObjects(e.LeftOperand, e.Operator, e.RightOperand)
}
