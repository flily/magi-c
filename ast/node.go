package ast

import (
	"reflect"

	"github.com/flily/magi-c/context"
)

type Comparable interface {
	context.ContextProvider
	EqualTo(archor context.ContextProvider, other Comparable) error
}

type Node interface {
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

func CheckNodeEqual[T Comparable](a T, b Comparable) (T, error) {
	va := reflect.ValueOf(a)
	switch va.Kind() {
	case reflect.Invalid:
		panic("CheckNodeEqual: `a` MUST NOT be untyped nil")

	case reflect.Pointer:
		vb := reflect.ValueOf(b)
		if vb.Kind() == reflect.Invalid {
			panic("CheckNodeEqual: `b` MUST NOT be untyped nil")
		}

		if vb.Kind() != reflect.Pointer {
			panic("CheckNodeEqual: only pointer type parameters are supported")
		}

		cb, ok := vb.Interface().(T)
		if !ok {
			return *new(T), a.Context().Error("expect a %T, got a %T", b, a).With("%T", b)
		}

		if va.IsNil() {
			if !vb.IsNil() {
				panic("CheckNodeEqual: `a` MUST NOT be typed nil when `b` is not nil")
			}

			return cb, nil
		}

		if vb.IsNil() {
			return *new(T), a.Context().Error("unexpected syntax element")
		}

		return cb, nil

	default:
		panic("CheckNodeEqual: only pointer kinds are supported")
	}
}

func CheckNilPointerEqual[T Comparable](archor context.ContextProvider, a T, b T) error {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Kind() != reflect.Pointer || vb.Kind() != reflect.Pointer {
		panic("CheckNilPointerEqual: only pointer kinds are supported")
	}

	if va.IsNil() && vb.IsNil() {
		return nil
	}

	if va.IsNil() {
		return archor.Context().NextInLineContext().Error("expect %T, got %T", b, a).With("%T", a)
	}

	if vb.IsNil() {
		return a.Context().Error("unexpected %T found", a).With("unexpected token")
	}

	return a.EqualTo(archor, b)
}

func CheckArrayEqual[T Comparable](message string, archor context.ContextProvider, a []T, b []T) error {
	if len(a) != len(b) {
		var ctx *context.Context
		if len(a) == 0 {
			ctx = archor.Context()

		} else {
			ctxList := make([]context.ContextProvider, 0, len(a))
			for _, item := range a {
				ctxList = append(ctxList, item)
			}
			ctx = context.JoinObjects(ctxList...)
		}
		return ctx.Error("wrong number of %s: expected %d, got %d", message, len(b), len(a))
	}

	for i, itemA := range a {
		if err := itemA.EqualTo(archor, b[i]); err != nil {
			return err
		}
	}

	return nil
}
