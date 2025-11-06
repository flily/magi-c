package parser

import (
	"go/ast"

	"github.com/flily/magi-c/context"
)

type Parser struct {
	Cursor *context.Cursor
}

func NewParser(filename string) (*Parser, error) {
	file, err := context.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	p := &Parser{
		Cursor: context.NewCursor(file),
	}

	return p, nil
}

func (p *Parser) Parse() (ast.Node, error) {
	return nil, nil
}
