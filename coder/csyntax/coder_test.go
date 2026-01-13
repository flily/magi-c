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

func TestCodeContext(t *testing.T) {
	ctx := NewContext("file.c", 10)

	checkInterfaceCodeElement(ctx)

	builder, writer := makeTestWriter(testStyle1)
	err := ctx.Write(writer, 0)
	if err != nil {
		t.Fatalf("CodeContext Write failed: %s", err)
	}

	expected := "#line 10 \"file.c\"\n"
	result := builder.String()
	if result != expected {
		t.Fatalf("CodeContext Write result wrong, expected '%s', got '%s'", expected, result)
	}
}

func TestElementCollectionBasic(t *testing.T) {
	collection := NewElementCollection(
		StringElement("lorem"), NewIntegerLiteral(42),
	)

	checkInterfaceCodeElement(collection)

	builder, writer := makeTestWriter(testStyle1)
	err := collection.Write(writer, 0)
	if err != nil {
		t.Fatalf("ElementCollection Write failed: %s", err)
	}

	expected := "lorem42"
	result := builder.String()
	if result != expected {
		t.Fatalf("ElementCollection Write result wrong, expected '%s', got '%s'", expected, result)
	}
}
