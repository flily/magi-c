package csyntax

import (
	"testing"

	"strings"
)

var (
	testStyle1 = KRStyle
	testStyle2 = &CodeStyle{
		Indent:                 "    ",
		FunctionBraceOnNewLine: true,
		FunctionBraceIndent:    "",
		IfBraceOnNewLine:       false,
		IfBraceIndent:          "",
		ForBraceOnNewLine:      false,
		ForBraceIndent:         "",
		WhileBraceOnNewLine:    false,
		WhileBraceIndent:       "",
		SwitchBraceOnNewLine:   false,
		SwitchBraceIndent:      "",
		CaseBranchIndent:       "",
		AssignmentSpacing:      true,
		BinaryOperationSpacing: true,
		TypeCastSpacing:        true,
		CommaSpacingBefore:     false,
		CommaSpacingAfter:      false,
		PointerSpacingBefore:   true,
		PointerSpacingAfter:    false,
	}
)

func makeTestWriter(sylte *CodeStyle) (*strings.Builder, *StyleWriter) {
	var b strings.Builder
	writer := sylte.MakeWriter(&b)
	return &b, writer
}

func TestCodeStyleClone(t *testing.T) {
	newStyle := KRStyle.Clone()

	newStyle.Indent = "\t"

	if KRStyle.Indent == newStyle.Indent {
		t.Fatalf("CodeStyle Clone failed: modifying clone affected original")
	}
}
