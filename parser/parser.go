package parser

import (
	"github.com/flily/magi-c/context"
)

type Parser struct {
	Filename          string
	Lines             []context.LineContent
	LineContexts      []*context.LineContext
	PreviousLineCount int

	line    int
	column  int
	pLine   int
	pColumn int
}

func NewParser(filename string) *Parser {
	p := &Parser{
		Filename:          filename,
		line:              0,
		column:            0,
		pLine:             0,
		pColumn:           0,
		PreviousLineCount: 5,
	}

	return p
}
