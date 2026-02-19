package coder

import (
	"testing"
)

func TestVariableMapBasic(t *testing.T) {
	vm := NewVariableMap()

	names := []string{"var1", "var2", "var3"}
	for _, name := range names {
		vm.Add(name)
	}

	if len(vm.Variables) != len(names) {
		t.Errorf("expect %d variables, got %d", len(names), len(vm.Variables))
	}

	for _, name := range names {
		if _, found := vm.Variables[name]; !found {
			t.Fatalf("variable '%s' not found", name)
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

	r := root.AddName("var1")
	if !r {
		t.Fatalf("failed to add variable 'var1' to root frame")
	}

	if _, found := root.GetName("var1"); !found {
		t.Fatalf("variable 'var1' should be found in root frame")
	}

	r = root.AddName("var1")
	if r {
		t.Fatalf("should not be able to add duplicate variable 'var1' to root frame")
	}

	child := NewFrameOn(root)
	if child.IsRoot() {
		t.Fatalf("child frame should not be root")
	}

	r = child.AddName("var2")
	if !r {
		t.Fatalf("failed to add variable 'var2' to child frame")
	}

	if _, found := child.GetName("var2"); !found {
		t.Fatalf("variable 'var2' should be found in child frame")
	}

	if _, found := child.GetName("var1"); found {
		t.Fatalf("variable 'var1' should not be found in child frame")
	}

	r = child.AddName("var1")
	if !r {
		t.Fatalf("should be able to add variable 'var1' to child frame even if it exists in parent frame")
	}

	if _, found := child.GetName("var1"); !found {
		t.Fatalf("variable 'var1' should be found in child frame")
	}

	if _, found := child.GetName("var2"); !found {
		t.Fatalf("variable 'var2' should be found in child frame")
	}
}

func TestContextGlobalUse(t *testing.T) {
	ctx := NewContext()
	r := ctx.IsGlobalContext()
	if !r {
		t.Fatalf("new context should be global context")
	}

	r = ctx.RegisterVariable("gvar1")
	if !r {
		t.Fatalf("failed to register variable 'gvar1' in global context")
	}

	if _, found := ctx.Find("gvar1"); !found {
		t.Fatalf("variable 'gvar1' should be found in global context")
	}

	r = ctx.RegisterVariable("gvar1")
	if r {
		t.Fatalf("should not be able to register duplicate variable 'gvar1' in global context")
	}

	r = ctx.RegisterVariable("gvar2")
	if !r {
		t.Fatalf("failed to register variable 'gvar2' in global context")
	}

	if _, found := ctx.Find("gvar1"); !found {
		t.Fatalf("variable 'gvar1' should still be found in global context")
	}

	if _, found := ctx.Find("gvar2"); !found {
		t.Fatalf("variable 'gvar2' should be found in global context")
	}
}

func TestContextFunctionUse(t *testing.T) {
	ctx := NewContext()
	ctx.RegisterVariable("gvar1")

	ctx.PushFrame()
	ctx.RegisterVariable("fvar1")

	if _, found := ctx.Find("fvar1"); !found {
		t.Fatalf("variable 'fvar1' should be found in function context")
	}

	if _, found := ctx.Find("gvar1"); !found {
		t.Fatalf("variable 'gvar1' should be found in function context")
	}

	ctx.PopFrame()
	if _, found := ctx.Find("fvar1"); found {
		t.Fatalf("variable 'fvar1' should not be found after popping function frame")
	}

	if _, found := ctx.Find("gvar1"); !found {
		t.Fatalf("variable 'gvar1' should still be found in global context")
	}
}
