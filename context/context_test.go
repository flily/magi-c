package context

import (
	"testing"

	"strings"
)

func TestLineContextHighlightText(t *testing.T) {
	//    0         1         2         3         4
	//    0    5    0    5    0    5    0    5    0
	s := "the quick brown fox jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s))
	lctx := line.MarkLine(4, 9)

	got := lctx.HighlighText("lorem ipsum")
	expected := strings.Join([]string{
		"  42:   the quick brown fox jumps over the lazy dog",
		"            ^^^^^",
		"            lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestLineContextHighlightTextMultiParts(t *testing.T) {
	//    0         1         2         3         4
	//    0    5    0    5    0    5    0    5    0
	s := "the quick brown fox jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s))
	lctx := line.MarkLine(4, 9)
	lctx.Mark(16, 19)

	got := lctx.HighlighText("lorem ipsum")
	expected := strings.Join([]string{
		"  42:   the quick brown fox jumps over the lazy dog",
		"            ^^^^^       ^^^",
		"            lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestLineContextHighlightTextWithTabCharacters(t *testing.T) {
	//    0         1         2         3         4
	//    0    5    0    5    0    5    0    5    0
	s := "the quick\tbrown fox jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s))
	lctx := line.MarkLine(10, 15)

	got := lctx.HighlighText("lorem ipsum")
	expected := strings.Join([]string{
		"  42:   the quick\tbrown fox jumps over the lazy dog",
		// 42:   the quick       brown fox jumps over the lazy dog
		"                        ^^^^^",
		"                        lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestLineContextHighlightTextWithTabCharactersMultiParts(t *testing.T) {
	//    0         1         2         3         4
	//    0    5    0    5    0    5    0    5    0
	s := "the quick\tbrown fox jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s))
	lctx := line.MarkLine(10, 15)
	lctx.Mark(4, 9)

	got := lctx.HighlighText("lorem ipsum")
	expected := strings.Join([]string{
		"  42:   the quick\tbrown fox jumps over the lazy dog",
		// 42:   the quick       brown fox jumps over the lazy dog
		"            ^^^^^       ^^^^^",
		"            lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestLineContextHighlightTextWithMixedLanguages1(t *testing.T) {
	//    0         1         2         3         4         5         6
	//    0    5    0    5    0    5    0    5    0    5    0    5    0
	s := "the quick brown fox 我能吞下玻璃而不伤身体 jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s))
	lctx := line.MarkLine(20, 24)

	got := lctx.HighlighText("lorem ipsum")
	expected := strings.Join([]string{
		"  42:   the quick brown fox 我能吞下玻璃而不伤身体 jumps over the lazy dog",
		"                            ^^^^^^^^",
		"                            lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestLineContextHighlightTextWithMixedLanguages2(t *testing.T) {
	//    0         1         2         3         4         5         6
	//    0    5    0    5    0    5    0    5    0    5    0    5    0
	s := "the quick brown fox 我能吞下玻璃而不伤身体 jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s))
	lctx := line.MarkLine(32, 37)

	got := lctx.HighlighText("lorem ipsum")
	expected := strings.Join([]string{
		"  42:   the quick brown fox 我能吞下玻璃而不伤身体 jumps over the lazy dog",
		"                                                   ^^^^^",
		"                                                   lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}
