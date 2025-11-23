package parser

import (
	"testing"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/preprocessor"
)

func runBasicTestOnCode(t *testing.T, code string) *ast.Document {
	t.Helper()

	parser := NewLLParserFromCode(code, "test.mc")
	preprocessor.RegisterPreprocessors(parser)
	document, err := parser.Parse()

	if err != nil {
		t.Fatalf("parse code failed:\n%s", err.Error())
	}

	if document == nil {
		t.Fatalf("document is nil")
	}

	return document
}
