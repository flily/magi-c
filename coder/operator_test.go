package coder

import (
	"testing"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/coder/csyntax"
)

func TestOperatorMap(t *testing.T) {
	op := ast.Plus
	expected := csyntax.OperatorAdd

	result := OperatorMap(op)

	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

func TestOperatorMapUnsupported(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for unsupported operator, but did not panic")
		}
	}()

	op := ast.Invalid // Assuming Invalid is not in the magicOperatorMap
	OperatorMap(op)
}
