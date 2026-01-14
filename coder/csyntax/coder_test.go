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
	elem.codeElement()
}

func checkInterfaceDeclaration(elem Declaration) {
	elem.declarationNode()
}

func checkInterfaceStatement(elem Statement) {
	elem.statementNode()
}

func checkInterfaceExpression(elem Expression) {
	elem.expressionNode()
}

func checkOutputResult(t *testing.T, builder *strings.Builder, expected string) {
	t.Helper()

	result := builder.String()
	if result != expected {
		t.Fatalf("Write result wrong:\nexpect:\n%s\ngot:\n%s", expected, result)
	}
}

func checkOutputOnStyle(t *testing.T, style *CodeStyle, expected string, elems ...CodeElement) {
	t.Helper()

	builder, writer := makeTestWriter(style)
	err := writer.Write(0, elems...)
	if err != nil {
		t.Fatalf("CodeElement write failed: %s", err)
	}

	checkOutputResult(t, builder, expected)
}

func TestCodeStyleClone(t *testing.T) {
	newStyle := KRStyle.Clone()

	newStyle.Indent = "\t"

	if KRStyle.Indent == newStyle.Indent {
		t.Fatalf("CodeStyle Clone failed: modifying clone affected original")
	}
}

func TestStyleWriterWriteStrings(t *testing.T) {
	expected := "helloworld"
	checkOutputOnStyle(t, testStyle1, expected,
		StringElement("hello"), StringElement("world"))
}

func TestStyleWriterWriteStringsWithDelimiter(t *testing.T) {
	expected := "hello world"
	checkOutputOnStyle(t, testStyle1, expected,
		StringElement("hello"), NewDelimiter(" "), StringElement("world"))
}

func TestStyleWriterWriteStringsWithDuplicatedDelimiters(t *testing.T) {
	expected := "hello world"
	checkOutputOnStyle(t, testStyle1, expected,
		StringElement("hello"), NewDelimiter(" "), NewDelimiter(" "), StringElement("world"))
}

func TestCodeContext(t *testing.T) {
	ctx := NewContext("file.c", 10)
	checkInterfaceCodeElement(ctx)

	expected := "#line 10 \"file.c\"\n"
	checkOutputOnStyle(t, testStyle1, expected, ctx)
}

func TestElementCollectionBasic(t *testing.T) {
	collection := NewElementCollection(
		StringElement("lorem"), NewIntegerLiteral(42),
	)

	checkInterfaceCodeElement(collection)

	expected := "lorem42"
	checkOutputOnStyle(t, testStyle1, expected, collection)
}

func TestElementCollectionSelect(t *testing.T) {
	c1 := NewElementCollection(StringElement("lorem"))
	c2 := NewElementCollection(StringElement("ipsum"))

	r1 := c1.Select(true, c2)
	r2 := c1.Select(false, c2)

	expected1 := "lorem"
	expected2 := "ipsum"
	checkOutputOnStyle(t, testStyle1, expected1, r1)
	checkOutputOnStyle(t, testStyle1, expected2, r2)
}
