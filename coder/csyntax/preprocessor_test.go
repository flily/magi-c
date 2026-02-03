package csyntax

import (
	"testing"

	"strings"
)

func TestPreprocessorIncludeWrite(t *testing.T) {
	ctx := makeLineContext("test.mc", 42)
	include := NewIncludeAngle(ctx, "stdio.h")

	checkInterfaceCodeElement(include)
	checkInterfaceDeclaration(include)
	checkInterfaceStatement(include)

	expected := strings.Join([]string{
		`#line 42 "test.mc"`,
		"#include <stdio.h>",
	}, "\n") + "\n"
	checkOutputOnStyle(t, KRStyle, expected, include)
}

func TestPreprocessorIncludeQuoteWrite(t *testing.T) {
	ctx := makeLineContext("test.mc", 100)
	include := NewIncludeQuote(ctx, "myheader.h")

	expected := strings.Join([]string{
		`#line 100 "test.mc"`,
		`#include "myheader.h"`,
	}, "\n") + "\n"
	checkOutputOnStyle(t, KRStyle, expected, include)
}
