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
