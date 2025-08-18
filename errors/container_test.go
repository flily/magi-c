package errors

import (
	"testing"

	"fmt"
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

		if e1.Lower() != base {
			t.Errorf("wrong lower error: %v", e1.Lower())
		}

		if e1.Inner() != nil {
			t.Errorf("inner error should be nil, got: %v", e1.Inner())
		}

		if e1.Message() != "lorem ipsum" {
			t.Errorf("message should be '%s', got: %s", "lorem ipsum", e1.Message())
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

func TestErrorMessageWithInnerError(t *testing.T) {
	base := New("base")

	err := fmt.Errorf("lorem ipsum")

	e := base.Derive("dolor sit amet").With(err)
	exp := "dolor sit amet < base [with: lorem ipsum]"
	got := e.Error()
	if got != exp {
		t.Errorf("wrong error message: %s", got)
		t.Errorf("           expected: %s", exp)
	}

	if e.Lower() != base {
		t.Errorf("wrong lower error: %v", e.Lower())
	}

	if e.Inner() != err {
		t.Errorf("inner error should be '%s', got: %v", "lorem ipsum", e.Inner())
	}
}
