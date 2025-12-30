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

	_, err := NewCoderFromBinary(source, "example.mc")
	if err != nil {
		t.Fatalf("NewCoderFromBinary failed:\n%s", err)
	}
}
