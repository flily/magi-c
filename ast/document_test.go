package ast

import (
	"testing"

	"strings"

	"github.com/flily/magi-c/context"
)

func generateTestWords(s string) []*context.Context {
	result := make([]*context.Context, 0, 100)

	cursor := context.NewCursorFromString("test.txt", s)
	for {
		cursor.SkipWhitespace()
		_, _, eof := cursor.Rune()
		if eof {
			break
		}

		i := 0
		start := cursor.State()
		for {
			c, eol, eof := cursor.Peek(i)
			if eof || eol {
				break
			}

			if c == ' ' || c == '\t' {
				break
			}

			i++
		}

		finish := cursor.PeekState(i)
		cursor.SetState(finish)

		_, ctx := cursor.FinishWith(start, finish)
		result = append(result, ctx)
	}

	return result
}

func TestBasicWordGenerate(t *testing.T) {
	text := strings.Join([]string{
		"aaa bbb ccc",
		"ddd",
	}, "\n")

	words := generateTestWords(text)
	expectedWords := []string{
		"aaa",
		"bbb",
		"ccc",
		"ddd",
	}

	if len(words) != len(expectedWords) {
		t.Fatalf("expected %d words, got %d", len(expectedWords), len(words))
	}

	for i, expected := range expectedWords {
		word := words[i].Content()
		if word != expected {
			t.Errorf("expected word %q, got %q", expected, word)
		}
	}
}

func TestDocument(t *testing.T) {
	text := strings.Join([]string{
		"# include < stdio.h >",
	}, "\n")

	ctxList := generateTestWords(text)

	decl := NewPreprocessorInclude(ctxList[0], ctxList[1],
		ctxList[2], ctxList[3], ctxList[4])

	doc := NewDocument(nil)
	doc.Add(decl)

	var _ Node = doc
	var _ Statement = doc
	doc.statementNode()

	if doc.Terminal() {
		t.Errorf("Document should not be terminal")
	}

	expected := ASTBuildDocument(
		ASTBuildIncludeAngle("stdio.h"),
	)

	if err := doc.EqualTo(nil, expected); err != nil {
		t.Errorf("Document not equal:\n%s", err)
	}
}

func TestDocumentNotEqualOnType(t *testing.T) {
	text := strings.Join([]string{
		"# include < stdio.h >",
	}, "\n")

	ctxList := generateTestWords(text)

	decl := NewPreprocessorInclude(ctxList[0], ctxList[1],
		ctxList[2], ctxList[3], ctxList[4])

	doc := NewDocument([]Declaration{decl})
	expected := ASTBuildValue(42)
	message := strings.Join([]string{
		"   1:   # include < stdio.h >",
		"        ^ ^^^^^^^ ^ ^^^^^^^ ^",
		"        expect a *ast.IntegerLiteral",
	}, "\n")
	err := doc.EqualTo(nil, expected)
	if err == nil {
		t.Fatalf("Document expected not equal, but equal")
	}
	if err.Error() != message {
		t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", message, err.Error())
	}
}
