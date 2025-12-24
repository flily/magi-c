package csyntax

import (
	"testing"

	"strings"
)

func TestPreprocessorIncludeWrite(t *testing.T) {
	ctx := NewContext("test.mc", 42)
	include := NewIncludeAngle(ctx, "stdio.h")

	var _ Declaration = include
	include.declarationNode()
	var _ Statement = include
	include.statementNode()

	builder, writer := makeTestWriter(KRStyle)
	err := include.Write(writer, 0)
	if err != nil {
		t.Fatalf("IncludeDirective Write failed: %s", err)
	}

	expected := strings.Join([]string{
		`#line 42 "test.mc"`,
		"#include <stdio.h>",
	}, "\n") + "\n"

	result := builder.String()
	if result != expected {
		t.Fatalf("IncludeDirective Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}

func TestPreprocessorIncludeQuoteWrite(t *testing.T) {
	ctx := NewContext("test.mc", 100)
	include := NewIncludeQuote(ctx, "myheader.h")

	builder, writer := makeTestWriter(KRStyle)
	err := include.Write(writer, 0)
	if err != nil {
		t.Fatalf("IncludeDirective Write failed: %s", err)
	}

	expected := strings.Join([]string{
		`#line 100 "test.mc"`,
		`#include "myheader.h"`,
	}, "\n") + "\n"

	result := builder.String()
	if result != expected {
		t.Fatalf("IncludeDirective Write result wrong:\nexpected:\n%s\ngot:\n%s", expected, result)
	}
}
