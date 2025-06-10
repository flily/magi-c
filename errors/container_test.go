package errors

import (
	"testing"
)

func TestErrorMessageWithBase(t *testing.T) {
	base := New("base")

	e1 := base.Derive("lorem ipsum")
	{
		expected := "lorem ipsum < base"
		if e1.Error() != expected {
			t.Errorf("wrong error message: %s", e1.Error())
			t.Errorf("           expected: %s", expected)
		}
	}

	e2 := e1.Derive("dolor sit amet")
	{
		expected := "dolor sit amet < lorem ipsum < base"
		if e2.Error() != expected {
			t.Errorf("wrong error message: %s", e2.Error())
			t.Errorf("           expected: %s", expected)
		}
	}
}
