package parser

import (
	"strings"
	"testing"

	"github.com/flily/magi-c/ast"
)

func TestLLParserSimpleStatement(t *testing.T) {
	code := strings.Join([]string{
		"#include <stdio.h>",
	}, "\n")

	document := runBasicTestOnCode(t, code)

	if len(document.Declarations) != 1 {
		t.Fatalf("Expected 1 declaration, got %d", len(document.Declarations))
	}

	include, ok := document.Declarations[0].(*ast.PreprocessorInclude)
	if !ok {
		t.Fatalf("Expected PreprocessorInclude, got %T", document.Declarations[0])
	}

	if value := include.Content.Content(); value != "stdio.h" {
		t.Fatalf("Expected include content 'stdio.h', got '%s'", value)
	}
}
