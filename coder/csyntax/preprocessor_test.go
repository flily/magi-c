package csyntax

import (
	"testing"

	"strings"
)

func TestPreprocessorIncludeWrite(t *testing.T) {
	ctx := makeLineContext("test.mc", 41)
	include := NewIncludeAngle(ctx, "stdio.h")

	checkInterfaceCodeElement(include)
	checkInterfaceDeclaration(include)
	checkInterfaceStatement(include)

	expected := strings.Join([]string{
		"#include <stdio.h>",
	}, "\n") + "\n"
	checkOutputOnStyle(t, KRStyle, expected, include)
}

func TestPreprocessorIncludeQuoteWrite(t *testing.T) {
	ctx := makeLineContext("test.mc", 99)
	include := NewIncludeQuote(ctx, "myheader.h")

	expected := strings.Join([]string{
		`#include "myheader.h"`,
	}, "\n") + "\n"
	checkOutputOnStyle(t, KRStyle, expected, include)
}

func TestInlineBlock(t *testing.T) {
	ctx := makeLineContext("test.mc", 7)
	inlineBlock := NewInlineBlock(ctx, "lorem ipsum;\ndolor sit amet;")

	checkInterfaceCodeElement(inlineBlock)
	checkInterfaceDeclaration(inlineBlock)
	checkInterfaceStatement(inlineBlock)

	expected := strings.Join([]string{
		`lorem ipsum;`,
		`dolor sit amet;`,
	}, "\n") + "\n"
	checkOutputOnStyle(t, KRStyle, expected, inlineBlock)
}
