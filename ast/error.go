package ast

import (
	"fmt"

	"github.com/flily/magi-c/context"
)

type Error struct {
	Message string
	Context *context.Context
}

func NewError(ctx *context.Context, message string, args ...any) *Error {
	e := &Error{
		Message: fmt.Sprintf(message, args...),
		Context: ctx,
	}

	return e
}

func (e *Error) Error() string {
	return e.Context.HighlightText(e.Message)
}
