package parser

import (
	"testing"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/preprocessor"
)

type llparserCorrectTestCase struct {
	Code     string
	Expected *ast.Document
}

func newCorrectCodeTestCase(code string, expected *ast.Document) *llparserCorrectTestCase {
	c := &llparserCorrectTestCase{
		Code:     code,
		Expected: expected,
	}

	return c
}

func (c *llparserCorrectTestCase) Run(t *testing.T) *ast.Document {
	t.Helper()

	parser := NewLLParserFromCode(c.Code, "test.mc")
	preprocessor.RegisterPreprocessors(parser)
	got, err := parser.Parse()

	if err != nil {
		t.Fatalf("parse code failed:\n%s", err.Error())
	}

	if got == nil {
		t.Fatalf("document is nil")
	}

	if err := got.EqualTo(nil, c.Expected); err != nil {
		t.Fatalf("expected document not equal to actual:\n%s", err)
	}

	return got
}
