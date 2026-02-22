package coder

import (
	"testing"
)

func TestVariableMapBasic(t *testing.T) {
	vm := NewVariableMap()

	names := []struct {
		nameInSource string
		nameInCode   string
	}{
		{"var1s", "var1c"},
		{"var2s", "var2c"},
		{"var3s", "var3c"},
	}
	for _, name := range names {
		vm.Add(name.nameInSource, name.nameInCode)
	}

	if len(vm.Variables) != len(names) {
		t.Errorf("expect %d variables, got %d", len(names), len(vm.Variables))
	}

	for _, name := range names {
		if got, found := vm.Get(name.nameInSource); !found || got.CodeName != name.nameInCode {
			t.Fatalf("variable '%s' not found or name mismatch: expect '%s', got '%s'", name.nameInSource, name.nameInCode, got.CodeName)
		}

	}

	_, found := vm.Get("nonexistent")
	if found {
		t.Fatalf("should not find nonexistent variable")
	}
}

func TestFrameBasic(t *testing.T) {
	root := NewFrame()
	if !root.IsRoot() {
		t.Fatalf("root frame should be root")
	}

	r := root.AddName("var1s", "var1c")
	if !r {
		t.Fatalf("failed to add variable 'var1' to root frame")
	}

	if got, found := root.GetName("var1s"); !found || got.CodeName != "var1c" {
		t.Fatalf("variable 'var1' not found or name mismatch: expect 'var1c', got '%s'", got.CodeName)
	}

	r = root.AddName("var1s", "var1c2")
	if r {
		t.Fatalf("should not be able to add duplicate variable 'var1' to root frame")
	}

	child := NewFrameOn(root)
	if child.IsRoot() {
		t.Fatalf("child frame should not be root")
	}

	r = child.AddName("var2s", "var2c")
	if !r {
		t.Fatalf("failed to add variable 'var2' to child frame")
	}

	if got, found := child.GetName("var2s"); !found || got.CodeName != "var2c" {
		t.Fatalf("variable 'var2' not found or name mismatch: expect 'var2c', got '%s'", got.CodeName)
	}

	if _, found := child.GetName("var1s"); found {
		t.Fatalf("variable 'var1' should not be found in child frame")
	}

	r = child.AddName("var1s", "var1c")
	if !r {
		t.Fatalf("should be able to add variable 'var1' to child frame even if it exists in parent frame")
	}

	if got, found := child.GetName("var1s"); !found || got.CodeName != "var1c" {
		t.Fatalf("variable 'var1' not found or name mismatch: expect 'var1c', got '%s'", got.CodeName)
	}

	if got, found := child.GetName("var2s"); !found || got.CodeName != "var2c" {
		t.Fatalf("variable 'var2' not found or name mismatch: expect 'var2c', got '%s'", got.CodeName)
	}
}

func TestContextGlobalUse(t *testing.T) {
	ctx := NewContext()
	r := ctx.IsGlobalContext()
	if !r {
		t.Fatalf("new context should be global context")
	}

	r = ctx.RegisterVariable("gvar1s", "gvar1c")
	if !r {
		t.Fatalf("failed to register variable 'gvar1' in global context")
	}

	if _, found := ctx.Find("gvar1s"); !found {
		t.Fatalf("variable 'gvar1' should be found in global context")
	}

	r = ctx.RegisterVariable("gvar1s", "gvar1c2")
	if r {
		t.Fatalf("should not be able to register duplicate variable 'gvar1' in global context")
	}

	r = ctx.RegisterVariable("gvar2s", "gvar2c")
	if !r {
		t.Fatalf("failed to register variable 'gvar2' in global context")
	}

	if _, found := ctx.Find("gvar1s"); !found {
		t.Fatalf("variable 'gvar1' should still be found in global context")
	}

	if _, found := ctx.Find("gvar2s"); !found {
		t.Fatalf("variable 'gvar2' should be found in global context")
	}
}

func TestContextFunctionUse(t *testing.T) {
	ctx := NewContext()
	ctx.RegisterVariable("gvar1s", "gvar1c")

	ctx.PushFrame()
	ctx.RegisterVariable("fvar1s", "fvar1c")

	if _, found := ctx.Find("fvar1s"); !found {
		t.Fatalf("variable 'fvar1' should be found in function context")
	}

	if _, found := ctx.Find("gvar1s"); !found {
		t.Fatalf("variable 'gvar1' should be found in function context")
	}

	ctx.PopFrame()
	if _, found := ctx.Find("fvar1s"); found {
		t.Fatalf("variable 'fvar1' should not be found after popping function frame")
	}

	if _, found := ctx.Find("gvar1s"); !found {
		t.Fatalf("variable 'gvar1' should still be found in global context")
	}
}
