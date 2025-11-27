package ast

import (
	"reflect"

	"github.com/flily/magi-c/context"
)

type Comparable interface {
	EqualTo(other Comparable) bool
}

type Node interface {
	context.ContextProvider
	Comparable
	HighlightText(message string, args ...any) string
}

type ContextContainer struct {
	context *context.Context
}

func (n *ContextContainer) Context() *context.Context {
	return n.context
}

func (n *ContextContainer) HighlightText(message string, args ...any) string {
	ctx := n.Context()
	return ctx.HighlightText(message, args...)
}

type TerminalNode interface {
	Node
	Type() TokenType
}

type TerminalNodeBase struct {
	ContextContainer
}

func NewTerminalNodeBase(ctx *context.Context) TerminalNodeBase {
	n := TerminalNodeBase{
		ContextContainer: ContextContainer{
			context: ctx,
		},
	}

	return n
}

type NonTerminalNode struct {
	provider context.ContextProvider
}

func (n *NonTerminalNode) Init(provider context.ContextProvider) {
	n.provider = provider
}

func (c *NonTerminalNode) HighlightText(message string, args ...any) string {
	ctx := c.provider.Context()
	return ctx.HighlightText(message, args...)
}

func CheckNodeEqual[T Comparable](a T, b Comparable) (T, bool) {
	va := reflect.ValueOf(a)
	switch va.Kind() {
	case reflect.Invalid:
		vb := reflect.ValueOf(b)
		return *new(T), vb.Kind() == reflect.Invalid

	case reflect.Pointer:
		vb := reflect.ValueOf(b)
		if vb.Kind() != reflect.Pointer {
			panic(" CheckNodeEqual: only pointer kinds are supported")
		}

		cb, ok := vb.Interface().(T)
		if !ok {
			return *new(T), false
		}

		if va.IsNil() && vb.IsNil() {
			return cb, true
		}

		return cb, true

	default:
		panic(" CheckNodeEqual: only pointer kinds are supported")
	}
}

func CheckNilPointerEqual[T Comparable](a T, b T) bool {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Kind() != reflect.Pointer || vb.Kind() != reflect.Pointer {
		panic(" CheckNilPointerEqual: only pointer kinds are supported")
	}

	if va.IsNil() && vb.IsNil() {
		return true
	}

	if va.IsNil() || vb.IsNil() {
		return false
	}

	return a.EqualTo(b)
}

func CheckArrayEqual[T Comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i, itemA := range a {
		if !itemA.EqualTo(b[i]) {
			return false
		}
	}

	return true
}
