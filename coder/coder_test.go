package coder

import (
	"testing"

	"strings"
)

func TestCoderFromBinary(t *testing.T) {
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
}
