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
		EOL:                    EOLLF,
	}
)

func makeTestWriter(sylte *CodeStyle) (*strings.Builder, *StyleWriter) {
	var b strings.Builder
	writer := sylte.MakeWriter(&b)
	return &b, writer
}

func checkInterfaceCodeElement(elem CodeElement) {
	var _ CodeElement = elem
	elem.codeElement()
}

func checkInterfaceDeclaration(elem Declaration) {
	var _ Declaration = elem
	elem.declarationNode()
}

func checkInterfaceStatement(elem Statement) {
	var _ Statement = elem
	elem.statementNode()
}

func checkInterfaceExpression(elem Expression) {
	var _ Expression = elem
	elem.expressionNode()
}

func TestCodeStyleClone(t *testing.T) {
	newStyle := KRStyle.Clone()

	newStyle.Indent = "\t"

	if KRStyle.Indent == newStyle.Indent {
		t.Fatalf("CodeStyle Clone failed: modifying clone affected original")
	}
}

func TestStyleWriterWriteStrings(t *testing.T) {
	builder, writer := makeTestWriter(testStyle1)

	err := writer.Write(0, StringElement("hello"), StringElement("world"))
	if err != nil {
		t.Fatalf("StyleWriter WriteStringItem failed: %s", err)
	}

	expected := "helloworld"
	result := builder.String()
	if result != expected {
		t.Fatalf("StyleWriter WriteStringItem result wrong, expected '%s', got '%s'", expected, result)
	}
}

func TestStyleWriterWriteStringsWithDelimiter(t *testing.T) {
	builder, writer := makeTestWriter(testStyle1)

	err := writer.Write(0, StringElement("hello"), NewDelimiter(" "), StringElement("world"))
	if err != nil {
		t.Fatalf("StyleWriter WriteStringItem failed: %s", err)
	}

	expected := "hello world"
	result := builder.String()
	if result != expected {
		t.Fatalf("StyleWriter WriteStringItem result wrong, expected '%s', got '%s'", expected, result)
	}
}

func TestStyleWriterWriteStringsWithDuplicatedDelimiters(t *testing.T) {
	builder, writer := makeTestWriter(testStyle1)

	err := writer.Write(0, StringElement("hello"), NewDelimiter(" "), NewDelimiter(" "), StringElement("world"))
	if err != nil {
		t.Fatalf("StyleWriter WriteStringItem failed: %s", err)
	}

	expected := "hello world"
	result := builder.String()
	if result != expected {
		t.Fatalf("StyleWriter WriteStringItem result wrong, expected '%s', got '%s'", expected, result)
	}
}
