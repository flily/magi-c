package parser

import (
	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/preprocessor"
)

type Parser interface {
	Parse() (*ast.Document, error)
	RegisterPreprocessor(name string, handler preprocessor.PreprocessorInitializer)
}

type Precedence int

const (
	PrecedenceLowest Precedence = iota
	PrecedenceLogicalOR
	PrecedenceLogicalAND
	PrecedenceComparisonEquality
	PrecedenceComparisonRelational
	PrecedenceSum
	PrecedenceProduct
	PrecedenceBitwiseOR
	PrecedenceBitwiseXOR
	PrecedenceBitwiseAND
	PrecedenceShift
	PrecedencePrefix
	PrecedenceCall
	PrecedenceIndex
	PrecedenceHighest
)

var precedenceMap = map[ast.TokenType]Precedence{
	ast.Plus: PrecedenceSum,
	ast.Sub:  PrecedenceSum,
}

func GetPrecedence(node ast.TerminalNode) Precedence {
	if node == nil {
		return PrecedenceLowest
	}

	t := node.Type()
	if prec, ok := precedenceMap[t]; ok {
		return prec
	}

	return PrecedenceLowest
}
