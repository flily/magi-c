package check

import (
	"testing"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
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

	conf := NewDefaultCheckConfigure()
	checker := NewCodeChecker(conf, doc)
	container := checker.Check()
	if container == nil {
		t.Fatalf("a DiagnosticContainer expected, but got nil")
	}

	if container.Count(context.Error) > 0 {
		t.Fatalf("code check expected to succeed, but got errors:\n%s", container.Error())
	}
}

func checkCodeError(t *testing.T, source string, expected string) {
	t.Helper()

	doc := parseCode(t, source)

	conf := NewDefaultCheckConfigure()
	checker := NewCodeChecker(conf, doc)
	container := checker.Check()
	if container == nil {
		t.Fatalf("code check expected to fail, but succeeded")
	}

	if container.Error() != expected {
		t.Fatalf("code check error mismatch, expected:\n%s\ngot:\n%s", expected, container.Error())
	}
}
