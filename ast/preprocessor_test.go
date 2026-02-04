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

	checkDeclarationNodeInterface(includeSegment)
	checkStatementNodeInterface(includeSegment)

	if includeSegment.Type() != NodePreprocessorInclude {
		t.Fatalf("include segment type expected %d, got %d", NodePreprocessorInclude, includeSegment.Type())
	}

	expected := ASTBuildIncludeAngle("stdio.h")
	if err := includeSegment.EqualTo(nil, expected); err != nil {
		t.Errorf("PreprocessorInclude not equal:\n%s", err)
	}
}

func TestPreprocessorIncludeNotEqualInContent(t *testing.T) {
	text := "# include < stdio.h >"
	ctxList := generateTestWords(text)

	includeSegment := NewPreprocessorInclude(ctxList[0], ctxList[1],
		ctxList[2], ctxList[3], ctxList[4])

	checkDeclarationNodeInterface(includeSegment)
	checkStatementNodeInterface(includeSegment)

	expected := ASTBuildIncludeAngle("stdlib.h")
	message := strings.Join([]string{
		"test.txt:1:13: error: wrong include content, expect 'stdlib.h', got 'stdio.h'",
		"    1 | # include < stdio.h >",
		"      |             ^^^^^^^",
		"      |             stdlib.h",
	}, "\n")

	err := includeSegment.EqualTo(nil, expected)
	if err == nil {
		t.Errorf("PreprocessorInclude expected not equal, but equal")
	}

	if err.Error() != message {
		t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
	}
}

func TestPreprocessorIncludeNotEqualInQuote(t *testing.T) {
	text := "# include < stdio.h >"
	ctxList := generateTestWords(text)

	includeSegment := NewPreprocessorInclude(ctxList[0], ctxList[1],
		ctxList[2], ctxList[3], ctxList[4])

	checkDeclarationNodeInterface(includeSegment)
	checkStatementNodeInterface(includeSegment)

	expected := ASTBuildIncludeQuote("stdio.h")
	message := strings.Join([]string{
		"test.txt:1:11: error: wrong include bracket, expect '\"' and '\"', got '<' and '>'",
		"    1 | # include < stdio.h >",
		"      |           ^         ^",
		`      |           "         "`,
	}, "\n")

	err := includeSegment.EqualTo(nil, expected)
	if err == nil {
		t.Errorf("PreprocessorInclude expected not equal, but equal")
	}

	if err.Error() != message {
		t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
	}
}

func TestPreprocessorIncludeNotEqual(t *testing.T) {
	text := "# include < stdio.h >"
	ctxList := generateTestWords(text)

	includeSegment := NewPreprocessorInclude(ctxList[0], ctxList[1],
		ctxList[2], ctxList[3], ctxList[4])

	checkDeclarationNodeInterface(includeSegment)
	checkStatementNodeInterface(includeSegment)

	expected := ASTBuildValue(42)
	message := strings.Join([]string{
		"test.txt:1:1: error: expect a *ast.IntegerLiteral, got a *ast.PreprocessorInclude",
		"    1 | # include < stdio.h >",
		"      | ^ ^^^^^^^ ^ ^^^^^^^ ^",
		"      | *ast.IntegerLiteral",
	}, "\n")

	err := includeSegment.EqualTo(nil, expected)
	if err == nil {
		t.Errorf("PreprocessorInclude expected not equal, but equal")
	}

	if err.Error() != message {
		t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
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

	checkDeclarationNodeInterface(inlineSegment)
	checkStatementNodeInterface(inlineSegment)

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

	checkDeclarationNodeInterface(inlineSegment)
	checkStatementNodeInterface(inlineSegment)

	if !inlineSegment.Empty() {
		t.Fatalf("expect inline segment is empty")
	}
}

func TestPreprocessorInlineNotEqualInType(t *testing.T) {
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

	checkDeclarationNodeInterface(inlineSegment)
	checkStatementNodeInterface(inlineSegment)

	expected := ASTBuildInline("c",
		"inline-content",
	)
	message := strings.Join([]string{
		"test.txt:1:10: error: wrong inline code type, expect 'c', got 'asm'",
		"    1 | # inline asm",
		"      |          ^^^",
		"      |          c",
	}, "\n")

	err := inlineSegment.EqualTo(nil, expected)
	if err == nil {
		t.Errorf("PreprocessorInline expected not equal, but equal")
	}

	if err.Error() != message {
		t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
	}
}

func TestPreprocessorInlineNotEqualInContent(t *testing.T) {
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

	checkDeclarationNodeInterface(inlineSegment)
	checkStatementNodeInterface(inlineSegment)

	expected := ASTBuildInline("asm",
		"different-content",
	)
	message := strings.Join([]string{
		"test.txt:2:1: error: wrong inline content",
		"    2 | inline-content",
		"      | ^^^^^^^^^^^^^^",
		"      | different-content",
	}, "\n")

	err := inlineSegment.EqualTo(nil, expected)
	if err == nil {
		t.Errorf("PreprocessorInline expected not equal, but equal")
	}

	if err.Error() != message {
		t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
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

	checkDeclarationNodeInterface(inlineSegment)
	checkStatementNodeInterface(inlineSegment)

	expected := ASTBuildValue(42)
	message := strings.Join([]string{
		"test.txt:1:1: error: expect a *ast.IntegerLiteral, got a *ast.PreprocessorInline",
		"    1 | # inline asm",
		"      | ^ ^^^^^^ ^^^",
		"    2 | inline-content",
		"      | ^^^^^^^^^^^^^^",
		"    3 | # end-inline asm",
		"      | ^ ^^^^^^^^^^ ^^^",
		"      | *ast.IntegerLiteral",
	}, "\n")

	err := inlineSegment.EqualTo(nil, expected)
	if err == nil {
		t.Errorf("PreprocessorInline expected not equal, but equal")
	}

	if err.Error() != message {
		t.Errorf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
	}
}
