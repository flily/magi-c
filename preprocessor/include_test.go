package preprocessor

import (
	"testing"

	"strings"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
)

func TestIncludeDirectiveAngleQuote(t *testing.T) {
	code := strings.Join([]string{
		"#include <stdio.h>",
	}, "\n")

	cursor := context.NewCursorFromString("example.c", code)
	node, err := scanDirectiveOn(cursor, Include)
	if err != nil {
		t.Fatalf("unexpected error:\n%v", err)
	}

	result, ok := node.(*ast.PreprocessorInclude)
	if !ok {
		t.Fatalf("expected PreprocessorInclude node, got %T", node)
	}

	gotHash := result.Hash.HighlightText("here")
	expHash := strings.Join([]string{
		"   1:   #include <stdio.h>",
		"        ^",
		"        here",
	}, "\n")
	if gotHash != expHash {
		t.Errorf("expected hash context highlight:\n%s\ngot:\n%s", expHash, gotHash)
	}

	gotCmd := result.Command.HighlightText("here")
	expCmd := strings.Join([]string{
		"   1:   #include <stdio.h>",
		"         ^^^^^^^",
		"         here",
	}, "\n")
	if gotCmd != expCmd {
		t.Errorf("expected command context highlight:\n%s\ngot:\n%s", expCmd, gotCmd)
	}

	gotLBracket := result.LBracket.HighlightText("here")
	expLBracket := strings.Join([]string{
		"   1:   #include <stdio.h>",
		"                 ^",
		"                 here",
	}, "\n")
	if gotLBracket != expLBracket {
		t.Errorf("expected LBracket context highlight:\n%s\ngot:\n%s", expLBracket, gotLBracket)
	}

	gotContent := result.Content.HighlightText("here")
	expContent := strings.Join([]string{
		"   1:   #include <stdio.h>",
		"                  ^^^^^^^",
		"                  here",
	}, "\n")
	if gotContent != expContent {
		t.Errorf("expected content context highlight:\n%s\ngot:\n%s", expContent, gotContent)
	}

	gotRBracket := result.RBracket.HighlightText("here")
	expRBracket := strings.Join([]string{
		"   1:   #include <stdio.h>",
		"                         ^",
		"                         here",
	}, "\n")
	if gotRBracket != expRBracket {
		t.Errorf("expected RBracket context highlight:\n%s\ngot:\n%s", expRBracket, gotRBracket)
	}
}
