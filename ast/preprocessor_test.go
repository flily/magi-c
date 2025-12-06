package ast

import (
	"strings"
	"testing"
)

func TestPreprocessorInclude(t *testing.T) {
	text := "# include < stdio.h >"
	ctxList := generateTestWords(text)

	includeSegment := NewPreprocessorInclude(ctxList[0], ctxList[1],
		ctxList[2], ctxList[3], ctxList[4])

	var _ Declaration = includeSegment
	includeSegment.declarationNode()
	var _ Statement = includeSegment
	includeSegment.statementNode()

	if includeSegment.Type() != NodePreprocessorInclude {
		t.Fatalf("include segment type expected %d, got %d", NodePreprocessorInclude, includeSegment.Type())
	}

	expected := ASTBuildIncludeAngle("stdio.h")
	if err := includeSegment.EqualTo(nil, expected); err != nil {
		t.Errorf("PreprocessorInclude not equal:\n%s", err)
	}
}

func TestPreprocessorIncludeNotEqual(t *testing.T) {
	text := "# include < stdio.h >"
	ctxList := generateTestWords(text)

	includeSegment := NewPreprocessorInclude(ctxList[0], ctxList[1],
		ctxList[2], ctxList[3], ctxList[4])

	{
		expected := ASTBuildIncludeAngle("stdlib.h")
		message := strings.Join([]string{
			"   1:   # include < stdio.h >",
			"                    ^^^^^^^",
			"                    wrong include content, expect 'stdlib.h', got 'stdio.h'",
		}, "\n")

		err := includeSegment.EqualTo(nil, expected)
		if err == nil {
			t.Errorf("PreprocessorInclude expected not equal, but equal")
		}

		if err.Error() != message {
			t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
		}
	}

	{
		expected := ASTBuildIncludeQuote("stdio.h")
		message := strings.Join([]string{
			"   1:   # include < stdio.h >",
			"                  ^         ^",
			"                  wrong include bracket, expect '\"' and '\"', got '<' and '>'",
		}, "\n")

		err := includeSegment.EqualTo(nil, expected)
		if err == nil {
			t.Errorf("PreprocessorInclude expected not equal, but equal")
		}

		if err.Error() != message {
			t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
		}
	}

	{
		expected := ASTBuildValue(42)
		message := strings.Join([]string{
			"   1:   # include < stdio.h >",
			"        ^ ^^^^^^^ ^ ^^^^^^^ ^",
			"        expect a *ast.IntegerLiteral",
		}, "\n")

		err := includeSegment.EqualTo(nil, expected)
		if err == nil {
			t.Errorf("PreprocessorInclude expected not equal, but equal")
		}

		if err.Error() != message {
			t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
		}
	}
}
