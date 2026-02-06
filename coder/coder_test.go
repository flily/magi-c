package coder

import (
	"testing"

	"strings"
)

func TestCoderFromBinary1(t *testing.T) {
	source := []byte(strings.Join([]string{
		`#include <stdio.h>`,
		`fun main() {`,
		`    #inline c`,
		`        printf("hello, world\n");`,
		`    #end-inline c`,
		`}`,
		``,
	}, "\n"))

	coder := NewCoder(".", ".")

	err := coder.ParseFileContent("example.mc", source)
	if err != nil {
		t.Fatalf("ParseFileContent failed:\n%s", err)
	}

	mainFilename := coder.FindMain()
	if mainFilename != "example.mc" {
		t.Fatalf("FindMain failed, expect 'example.mc', got '%s'", mainFilename)
	}
}

func TestCoderFromBinary2(t *testing.T) {
	source := []byte(strings.Join([]string{
		`#include <stdio.h>`,
		`fun helper() {`,
		`    #inline c`,
		`        printf("helper function\n");`,
		`    #end-inline c`,
		`}`,
		``,
	}, "\n"))

	coder := NewCoder(".", ".")

	err := coder.ParseFileContent("example.mc", source)
	if err != nil {
		t.Fatalf("ParseFileContent failed:\n%s", err)
	}

	mainFilename := coder.FindMain()
	if mainFilename != "" {
		t.Fatalf("FindMain failed, expect empty string, got '%s'", mainFilename)
	}
}

func TestOutputFilename(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"example.mc", "example.mc.c"},
		{"src/main.mc", "src/main.mc.c"},
		{"test/test1.mc", "test/test1.mc.c"},
	}

	for _, c := range cases {
		output := OutputFilename(c.input)
		if output != c.expected {
			t.Fatalf("OutputFilename failed for input '%s', expected '%s', got '%s'", c.input, c.expected, output)
		}
	}
}
