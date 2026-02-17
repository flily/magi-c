package coder

import (
	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/coder/csyntax"
)

var magicOperatorMap = map[ast.TokenType]csyntax.Punctuator{
	ast.Plus: csyntax.OperatorAdd,
	ast.Sub:  csyntax.OperatorSubtract,
}

func OperatorMap(op ast.TokenType) csyntax.Punctuator {
	p, found := magicOperatorMap[op]
	if found {
		return p
	}

	panic("unsupported operator: " + op.String())
}
