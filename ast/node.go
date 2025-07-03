package ast

import (
	"github.com/flily/magi-c/token"
)

type Node interface {
	Terminal() bool
	Context() *token.Context
}
