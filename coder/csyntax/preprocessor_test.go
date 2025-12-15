package csyntax

import (
	"strings"
	"testing"
)

func TestPreprocessorIncludeWrite(t *testing.T) {
	ctx := NewContext("test.mc", 42)
	include := NewIncludeAngle(ctx, "stdio.h")

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
