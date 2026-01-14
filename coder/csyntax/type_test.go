package csyntax

import (
	"testing"
)

func TestTypeWriteOnConcreteType(t *testing.T) {
	ty := NewConcreteType("int")

	checkInterfaceCodeElement(ty)

	expected := "int"
	checkOutputOnStyle(t, testStyle1, expected, ty)
}

func TestTypeWriteOnPointerType1(t *testing.T) {
	ty := NewPointerType("int")
	checkInterfaceCodeElement(ty)

	expected := "int* "
	checkOutputOnStyle(t, testStyle1, expected, ty)
}

func TestTypeWriteOnPointerType2(t *testing.T) {
	ty := NewType("char", 2)
	checkInterfaceCodeElement(ty)

	expected := "char **"
	checkOutputOnStyle(t, testStyle2, expected, ty)
}
