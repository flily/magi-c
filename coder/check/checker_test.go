package check

import (
	"testing"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/parser"
	"github.com/flily/magi-c/preprocessor"
)

func parseCode(t *testing.T, source string) *ast.Document {
	t.Helper()

	parser := parser.NewLLParserFromCode(source, "test.mc")
	preprocessor.RegisterPreprocessors(parser)
	doc, err := parser.Parse()
	if err != nil {
		t.Fatalf("parse code failed:\n%s", err.Error())
	}

	if doc == nil {
		t.Fatalf("document is nil")
	}

	return doc
}

func checkCodeCorrect(t *testing.T, source string) {
	t.Helper()

	doc := parseCode(t, source)

	checker := NewCodeChecker(doc)
	err := checker.Check()
	if err != nil {
		t.Fatalf("code check failed:\n%s", err.Error())
	}
}

func checkCodeError(t *testing.T, source string, expected string) {
	t.Helper()

	doc := parseCode(t, source)

	checker := NewCodeChecker(doc)
	err := checker.Check()
	if err == nil {
		t.Fatalf("code check expected to fail, but succeeded")
	}

	if err.Error() != expected {
		t.Fatalf("code check error mismatch, expected:\n%s\ngot:\n%s", expected, err.Error())
	}
}
