package csyntax

import (
	"testing"
)

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

func TestLevel(t *testing.T) {
	l := NewDefaultLevel()
	if l.IndentLevel != 0 || l.ParanthesisLevel != 0 {
		t.Fatalf("NewDefaultLevel failed: expected indent and paranthesis levels to be 0")
	}

	l = l.NextIndent()
	if l.IndentLevel != 1 || l.ParanthesisLevel != 0 {
		t.Fatalf("NextIndent failed: expected indent level to be 1 and paranthesis level to be 0")
	}

	l = l.NextParanthesis()
	if l.IndentLevel != 1 || l.ParanthesisLevel != 1 {
		t.Fatalf("NextParanthesis failed: expected indent level to be 1 and paranthesis level to be 1")
	}

	l = l.NextIndent()
	if l.IndentLevel != 2 || l.ParanthesisLevel != 0 {
		t.Fatalf("NextIndent failed: expected indent level to be 2 and paranthesis level to be reset to 0")
	}
}
