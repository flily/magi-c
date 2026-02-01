package ast

import (
	"testing"

	"strings"
)

func TestCheckNodeEqualAIsUntypedNil(t *testing.T) {
	exp := "CheckNodeEqual: `a` MUST NOT be untyped nil"
	ctxList := generateTestWords("lorem ipsum")
	c1 := NewStringLiteral(ctxList[0], "lorem")

	defer func() {
		r := recover()
		if r != exp {
			t.Errorf("expect: %s\n", exp)
			t.Errorf("   got: %v\n", r)
			t.Fatalf("wrong panic message")
		}
	}()

	_, _ = CheckNodeEqual[Comparable](nil, c1)
	t.Fatalf("CheckNodeEqual MUST panic")
}

func TestCheckNodeEqualBIsUntypedNil(t *testing.T) {
	exp := "CheckNodeEqual: `b` MUST NOT be untyped nil"
	ctxList := generateTestWords("lorem ipsum")
	c1 := NewStringLiteral(ctxList[0], "lorem")

	defer func() {
		r := recover()
		if r != exp {
			t.Errorf("expect: %s\n", exp)
			t.Errorf("   got: %v\n", r)
			t.Fatalf("wrong panic message")
		}
	}()

	_, _ = CheckNodeEqual[Comparable](c1, nil)
	t.Fatalf("CheckNodeEqual MUST panic")
}

func TestCheckNodeEqualABAreBothNil(t *testing.T) {
	var a *StringLiteral = nil
	var b *StringLiteral = nil

	_, err := CheckNodeEqual[Comparable](a, b)
	if err != nil {
		t.Fatalf("expected both nil nodes equal")
	}
}

func TestCheckNodeEqualAIsTypedNil(t *testing.T) {
	exp := "CheckNodeEqual: `a` MUST NOT be typed nil when `b` is not nil"
	ctxList := generateTestWords("lorem ipsum")
	var a *StringLiteral = nil
	c1 := NewStringLiteral(ctxList[0], "lorem")

	defer func() {
		r := recover()
		if r != exp {
			t.Errorf("expect: %s\n", exp)
			t.Errorf("   got: %v\n", r)
			t.Fatalf("wrong panic message")
		}
	}()

	_, _ = CheckNodeEqual(a, c1)
	t.Fatalf("CheckNodeEqual MUST panic")
}

func TestCheckNodeEqualBIsTypedNil(t *testing.T) {
	ctxList := generateTestWords("lorem ipsum")
	a := NewStringLiteral(ctxList[0], "lorem")
	var b *StringLiteral = nil

	expA := strings.Join([]string{
		"test.txt:1:1: error: unexpected syntax element",
		"    1 | lorem ipsum",
		"      | ^^^^^",
	}, "\n")
	_, err := CheckNodeEqual(a, b)
	if err == nil {
		t.Fatalf("expected one nil node not equal to the other")
	}

	if err.Error() != expA {
		t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", expA, err.Error())
	}
}
