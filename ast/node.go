package ast

import "github.com/flily/magi-c/context"

type Node interface {
	Terminal() bool
	Context() *context.Context
}
