package ast

import (
	"github.com/flily/magi-c/context"
)

// fun Name()
type Func struct {
	Function     *context.Context
	FunctionName *context.Context
	LeftParan    *context.Context
	RightParan   *context.Context
}

func (f *Func) Terminal() bool {
	return false
}

func (f *Func) Context() *context.Context {
	return context.Join(
		f.Function,
		f.FunctionName,
		f.LeftParan,
		f.RightParan,
	)
}
