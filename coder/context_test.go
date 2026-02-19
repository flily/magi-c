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
