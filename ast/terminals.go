package ast

import (
	"fmt"
	"reflect"

	"github.com/flily/magi-c/context"
)

type EOLLiteral struct {
	TerminalNodeBase
	EOL string
}

func NewEOLLiteral(ctx *context.Context, eol string) *EOLLiteral {
	l := &EOLLiteral{
		TerminalNodeBase: NewTerminalNodeBase(ctx),
		EOL:              eol,
	}

	return l
}

func (l *EOLLiteral) expressionNode() {}

func (l *EOLLiteral) Type() TokenType {
	return EOL
}

func (l *EOLLiteral) EqualTo(_ context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(l, other)
	if err != nil {
		return err
	}

	if o.EOL != l.EOL {
		return NewError(l.Context(), "wrong EOL value, expect %s, got %s", o.EOL, l.EOL)
	}

	return nil
}

type StringLiteral struct {
	TerminalNodeBase
	Value string
}

func NewStringLiteral(ctx *context.Context, value string) *StringLiteral {
	l := &StringLiteral{
		TerminalNodeBase: NewTerminalNodeBase(ctx),
		Value:            value,
	}

	return l
}

func (l *StringLiteral) expressionNode() {}

func (l *StringLiteral) Type() TokenType {
	return String
}

func (l *StringLiteral) EqualTo(_ context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(l, other)
	if err != nil {
		return err
	}

	if o.Value != l.Value {
		return NewError(l.Context(), "wrong string value, expect '%s', got '%s'", o.Value, l.Value)
	}

	return nil
}

type IntegerLiteral struct {
	TerminalNodeBase
	Value uint64
}

func NewIntegerLiteral(ctx *context.Context, value uint64) *IntegerLiteral {
	l := &IntegerLiteral{
		TerminalNodeBase: NewTerminalNodeBase(ctx),
		Value:            value,
	}

	return l
}

func (l *IntegerLiteral) expressionNode() {}

func (l *IntegerLiteral) Type() TokenType {
	return Integer
}

func (l *IntegerLiteral) EqualTo(_ context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(l, other)
	if err != nil {
		return err
	}

	if o.Value != l.Value {
		return NewError(l.Context(), "wrong integer value, expect %v, got %v", o.Value, l.Value)
	}

	return nil
}

type FloatLiteral struct {
	TerminalNodeBase
	Value float64
}

func NewFloatLiteral(ctx *context.Context, value float64) *FloatLiteral {
	l := &FloatLiteral{
		TerminalNodeBase: NewTerminalNodeBase(ctx),
		Value:            value,
	}

	return l
}

func (l *FloatLiteral) expressionNode() {}

func (l *FloatLiteral) Type() TokenType {
	return Float
}

func (l *FloatLiteral) EqualTo(_ context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(l, other)
	if err != nil {
		return err
	}

	if o.Value != l.Value {
		return NewError(l.Context(), "wrong float value, expect %v, got %v", o.Value, l.Value)
	}

	return nil
}

func ASTBuildValue(v any) Expression {
	switch val := v.(type) {
	case string:
		return NewStringLiteral(nil, val)

	case uint, uint8, uint16, uint32, uint64:
		value := reflect.ValueOf(val)
		return NewIntegerLiteral(nil, value.Uint())

	case int, int8, int16, int32, int64:
		value := reflect.ValueOf(val)
		return NewIntegerLiteral(nil, uint64(value.Int()))

	case float32, float64:
		return NewFloatLiteral(nil, val.(float64))

	default:
		s := fmt.Sprintf("ASTBuildValue: unsupported value type: %T", v)
		panic(s)
	}
}

type Identifier struct {
	TerminalNodeBase
	Name string
}

func NewIdentifier(ctx *context.Context, name string) *Identifier {
	id := &Identifier{
		TerminalNodeBase: NewTerminalNodeBase(ctx),
		Name:             name,
	}

	return id
}

func ASTBuildIdentifier(name string) *Identifier {
	return NewIdentifier(nil, name)
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) Type() TokenType {
	return IdentifierName
}

func (i *Identifier) EqualTo(_ context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(i, other)
	if err != nil {
		return err
	}

	if o.Name != i.Name {
		return NewError(i.Context(), "wrong identifier name, expect %s, got %s", o.Name, i.Name)
	}

	return nil
}

type TerminalToken struct {
	TerminalNodeBase
	Token TokenType
}

func NewTerminalToken(ctx *context.Context, token TokenType) *TerminalToken {
	t := &TerminalToken{
		TerminalNodeBase: NewTerminalNodeBase(ctx),
		Token:            token,
	}

	return t
}

func ASTBuildSymbol(token TokenType) *TerminalToken {
	return NewTerminalToken(nil, token)
}

func ASTBuildKeyword(token TokenType) *TerminalToken {
	return NewTerminalToken(nil, token)
}

func (t *TerminalToken) Type() TokenType {
	return t.Token
}

func (t *TerminalToken) EqualTo(_ context.ContextProvider, other Comparable) error {
	o, err := CheckNodeEqual(t, other)
	if err != nil {
		return err
	}

	if o.Token != t.Token {
		return NewError(t.Context(), "wrong token type, expect '%v', got '%v'", o.Token, t.Token)
	}

	return nil
}
