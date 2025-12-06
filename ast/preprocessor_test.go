package ast

import (
	"testing"

	"strings"
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

func TestPreprocessorInlineEqual(t *testing.T) {
	codeType := "asm"
	content := "inline-content"
	text := strings.Join([]string{
		"# inline asm",
		content,
		"# end-inline asm",
	}, "\n")
	ctxList := generateTestWords(text)

	inlineSegment := NewPreprocessorInline(ctxList[0], ctxList[1], codeType, ctxList[2],
		content, ctxList[3],
		ctxList[4], ctxList[5], ctxList[6])

	var _ Declaration = inlineSegment
	inlineSegment.declarationNode()
	var _ Statement = inlineSegment
	inlineSegment.statementNode()

	if inlineSegment.Type() != NodePreprocessorInline {
		t.Fatalf("inline segment type expected %d, got %d", NodePreprocessorInline, inlineSegment.Type())
	}

	expected := ASTBuildInline("asm",
		"inline-content",
	)
	if err := inlineSegment.EqualTo(nil, expected); err != nil {
		t.Errorf("PreprocessorInline not equal:\n%s", err)
	}
}

func TestProprocessorInlineIsEmpty(t *testing.T) {
	text := strings.Join([]string{
		"# inline c",
		"# end-inline c",
	}, "\n")
	ctxList := generateTestWords(text)

	inlineSegment := NewPreprocessorInline(ctxList[0], ctxList[1], "c", ctxList[2],
		"", nil,
		ctxList[3], ctxList[4], ctxList[5])

	if !inlineSegment.Empty() {
		t.Fatalf("expect inline segment is empty")
	}
}

func TestPreprocessorInlineNotEqual(t *testing.T) {
	codeType := "asm"
	content := "inline-content"
	text := strings.Join([]string{
		"# inline asm",
		content,
		"# end-inline asm",
	}, "\n")
	ctxList := generateTestWords(text)

	inlineSegment := NewPreprocessorInline(ctxList[0], ctxList[1], codeType, ctxList[2],
		content, ctxList[3],
		ctxList[4], ctxList[5], ctxList[6])

	{
		expected := ASTBuildInline("c",
			"inline-content",
		)
		message := strings.Join([]string{
			"   1:   # inline asm",
			"                 ^^^",
			"                 wrong inline code type, expect 'c', got 'asm'",
		}, "\n")

		err := inlineSegment.EqualTo(nil, expected)
		if err == nil {
			t.Errorf("PreprocessorInline expected not equal, but equal")
		}

		if err.Error() != message {
			t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
		}
	}

	{
		expected := ASTBuildInline("asm",
			"different-content",
		)
		message := strings.Join([]string{
			"   2:   inline-content",
			"        ^^^^^^^^^^^^^^",
			"        wrong inline content, expect 'different-content', got 'inline-content'",
		}, "\n")

		err := inlineSegment.EqualTo(nil, expected)
		if err == nil {
			t.Errorf("PreprocessorInline expected not equal, but equal")
		}

		if err.Error() != message {
			t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
		}
	}

	{
		expected := ASTBuildValue(42)
		message := strings.Join([]string{
			"   1:   # inline asm",
			"        ^ ^^^^^^ ^^^",
			"   2:   inline-content",
			"        ^^^^^^^^^^^^^^",
			"   3:   # end-inline asm",
			"        ^ ^^^^^^^^^^ ^^^",
			"        expect a *ast.IntegerLiteral",
		}, "\n")

		err := inlineSegment.EqualTo(nil, expected)
		if err == nil {
			t.Errorf("PreprocessorInline expected not equal, but equal")
		}

		if err.Error() != message {
			t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
		}
	}
}
