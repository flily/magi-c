package coder

import (
	"bytes"
	"testing"

	"strings"
)

const (
	testFilename = "test.mc"
)

func testOutputCode(t *testing.T, code string, expected string) {
	t.Helper()

	coder := NewCoder(".", ".")
	err := coder.ParseFileContent(testFilename, []byte(code))
	if err != nil {
		t.Fatalf("ParseFileContent failed:\n%s", err)
	}

	err = coder.Check(testFilename)
	if err != nil {
		t.Fatalf("Check failed:\n%s", err)
	}

	buf := bytes.NewBuffer(nil)
	err = coder.OutputTo(testFilename, buf)
	if err != nil {
		t.Fatalf("OutputTo failed:\n%s", err)
	}

	output := buf.String()
	if output != expected {
		t.Fatalf("Output code mismatch:\nExpect:\n%s\nGot:\n%s", expected, output)
	}
}

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

func TestCodeBasicGenerate(t *testing.T) {
	source := strings.Join([]string{
		`#include <stdio.h>`,
		`fun main() {`,
		`    #inline c`,
		`    printf("hello, world\n");`,
		`    #end-inline c`,
		`}`,
		``,
	}, "\n")

	expected := strings.Join([]string{
		`#line 1 "test.mc"`,
		`#include <stdio.h>`,
		`void main()`,
		`{`,
		`#line 3 "test.mc"`,
		`    printf("hello, world\n");`,
		`}`,
		``,
	}, "\n")

	testOutputCode(t, source, expected)
}
