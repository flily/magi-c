package csyntax

import (
	"testing"
)

func TestTypeWriteOnConcreteType(t *testing.T) {
	ty := NewConcreteType("int")
	b, writer := makeTestWriter(KRStyle)
	err := ty.Write(writer, 0)
	if err != nil {
		t.Fatalf("Type Write failed: %s", err)
	}

	expected := "int"
	result := b.String()
	if result != expected {
		t.Fatalf("Type Write result wrong:\nexpected: %s\ngot: %s", expected, result)
	}
}

func TestTypeWriteOnPointerType1(t *testing.T) {
	ty := NewPointerType("int")
	b, writer := makeTestWriter(KRStyle)
	err := ty.Write(writer, 0)
	if err != nil {
		t.Fatalf("Type Write failed: %s", err)
	}

	expected := "int* "
	result := b.String()
	if result != expected {
		t.Fatalf("Type Write result wrong:\nexpected: %s\ngot: %s", expected, result)
	}
}

func TestTypeWriteOnPointerType2(t *testing.T) {
	ty := NewType("char", 2)

	style := KRStyle.Clone()
	style.PointerSpacingBefore = true
	b, writer := makeTestWriter(style)
	err := ty.Write(writer, 0)
	if err != nil {
		t.Fatalf("Type Write failed: %s", err)
	}

	expected := "char ** "
	result := b.String()
	if result != expected {
		t.Fatalf("Type Write result wrong:\nexpected: %s\ngot: %s", expected, result)
	}
}
