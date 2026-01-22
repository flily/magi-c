package context

import (
	"fmt"
)

type Error struct {
	Message string
	Context *Context
}

func NewError(ctx *Context, message string, args ...any) *Error {
	e := &Error{
		Message: fmt.Sprintf(message, args...),
		Context: ctx,
	}

	return e
}

func (e *Error) Error() string {
	return e.Context.HighlightText(e.Message)
}
