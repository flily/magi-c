package preprocessor

import (
	"testing"

	"strings"

	"github.com/flily/magi-c/context"
)

func TestScanDirective(t *testing.T) {
	code := strings.Join([]string{
		"#include <stdio.h>",
	}, "\n")

	cursor := context.NewCursorFromString("example.c", code)
	cmd, hashCtx, cmdCtx, err := ScanDirective(cursor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cmd != "include" {
		t.Errorf("expected command 'include', got '%s'", cmd)
	}

	hashExp := strings.Join([]string{
		"   1:   #include <stdio.h>",
		"        ^",
		"        here",
	}, "\n")
	hashGot := hashCtx.HighlightText("here")
	if hashGot != hashExp {
		t.Errorf("expected hash context highlight:\n%s\ngot:\n%s", hashExp, hashGot)
	}

	cmdExp := strings.Join([]string{
		"   1:   #include <stdio.h>",
		"         ^^^^^^^",
		"         here",
	}, "\n")
	cmdGot := cmdCtx.HighlightText("here")
	if cmdGot != cmdExp {
		t.Errorf("expected command context highlight:\n%s\ngot:\n%s", cmdExp, cmdGot)
	}
}
