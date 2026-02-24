package csyntax

import (
	"testing"

	"strings"

	"github.com/flily/magi-c/context"
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

func makeLineContext(filename string, line int) *context.Context {
	fileContext := &context.FileContext{
		Filename: filename,
	}
	l := context.NewLineFromString(line, "lorem ipsum", "\n")

	lineContext := &context.LineContext{
		File:    fileContext,
		Content: &l,
	}

	return lineContext.Mark(0, 10)
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
	elem1 := StringElement("hello")
	elem2 := StringElement("world")

	checkInterfaceCodeElement(elem1)
	checkInterfaceCodeElement(elem2)

	expected := "helloworld"
	checkOutputOnStyle(t, testStyle1, expected, elem1, elem2)
}

func TestStyleWriterWriteStringsWithDelimiter(t *testing.T) {
	elem1 := StringElement("hello")
	delimiter := NewDelimiter(" ")
	elem2 := StringElement("world")

	checkInterfaceCodeElement(elem1)
	checkInterfaceCodeElement(delimiter)
	checkInterfaceCodeElement(elem2)

	expected := "hello world"
	checkOutputOnStyle(t, testStyle1, expected,
		elem1, delimiter, elem2)
}

func TestStyleWriterWriteStringsWithDuplicatedDelimiters(t *testing.T) {
	elem1 := StringElement("hello")
	delimiter1 := NewDelimiter(" ")
	delimiter2 := NewDelimiter(" ")
	elem2 := StringElement("world")

	checkInterfaceCodeElement(elem1)
	checkInterfaceCodeElement(delimiter1)
	checkInterfaceCodeElement(delimiter2)
	checkInterfaceCodeElement(elem2)

	expected := "hello world"
	checkOutputOnStyle(t, testStyle1, expected,
		elem1, delimiter1, delimiter2, elem2)
}

func TestCodeContext(t *testing.T) {
	ctx := makeLineContext("file.c", 9)

	cctx := NewContext(ctx)
	checkInterfaceCodeElement(cctx)
	checkInterfaceDeclaration(cctx)
	checkInterfaceStatement(cctx)

	expected := "#line 10 \"file.c\"\n"
	checkOutputOnStyle(t, testStyle1, expected, cctx)
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
