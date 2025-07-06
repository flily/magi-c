package ast

import (
	"github.com/flily/magi-c/token"
)

// fun Name()
type Func struct {
	Function     *token.Context
	FunctionName *token.Context
	LeftParan    *token.Context
	RightParan   *token.Context
}

func (f *Func) Terminal() bool {
	return false
}

func (f *Func) Context() *token.Context {
	return token.JoinContexts(f.Function,
		f.FunctionName,
		f.LeftParan,
		f.RightParan)
}
