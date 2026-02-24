package csyntax

import (
	"testing"

	"strings"
)

func TestPreprocessorIncludeWrite(t *testing.T) {
	include := NewIncludeAngle("stdio.h")

	checkInterfaceCodeElement(include)
	checkInterfaceDeclaration(include)
	checkInterfaceStatement(include)

	expected := strings.Join([]string{
		"#include <stdio.h>",
	}, "\n") + "\n"
	checkOutputOnStyle(t, KRStyle, expected, include)
}

func TestPreprocessorIncludeQuoteWrite(t *testing.T) {
	include := NewIncludeQuote("myheader.h")

	expected := strings.Join([]string{
		`#include "myheader.h"`,
	}, "\n") + "\n"
	checkOutputOnStyle(t, KRStyle, expected, include)
}

func TestInlineBlock(t *testing.T) {
	inlineBlock := NewInlineBlock("lorem ipsum;\ndolor sit amet;")

	checkInterfaceCodeElement(inlineBlock)
	checkInterfaceDeclaration(inlineBlock)
	checkInterfaceStatement(inlineBlock)

	expected := strings.Join([]string{
		`lorem ipsum;`,
		`dolor sit amet;`,
	}, "\n") + "\n"
	checkOutputOnStyle(t, KRStyle, expected, inlineBlock)
}
