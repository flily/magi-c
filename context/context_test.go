package context

import (
	"strings"
	"testing"
)

func TestLineContextHighlight(t *testing.T) {
	//    0         1         2         3         4
	//    0    5    0    5    0    5    0    5    0
	s := "the quick brown fox jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s))
	ctx := line.MarkLine(4, 9)

	got := ctx.Highlight("lorem ipsum")
	expected := strings.Join([]string{
		"  42:   the quick brown fox jumps over the lazy dog",
		"            ^^^^^",
		"            lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestLineContextHighlightMultiParts(t *testing.T) {
	//    0         1         2         3         4
	//    0    5    0    5    0    5    0    5    0
	s := "the quick brown fox jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s))
	ctx := line.MarkLine(4, 9)
	ctx.Mark(16, 19)

	got := ctx.Highlight("lorem ipsum")
	expected := strings.Join([]string{
		"  42:   the quick brown fox jumps over the lazy dog",
		"            ^^^^^       ^^^",
		"            lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestLineContextHighlightWithMixedLanguages1(t *testing.T) {
	//    0         1         2         3         4         5         6
	//    0    5    0    5    0    5    0    5    0    5    0    5    0
	s := "the quick brown fox 我能吞下玻璃而不伤身体 jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s))
	ctx := line.MarkLine(20, 24)

	got := ctx.Highlight("lorem ipsum")
	expected := strings.Join([]string{
		"  42:   the quick brown fox 我能吞下玻璃而不伤身体 jumps over the lazy dog",
		"                            ^^^^^^^^",
		"                            lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}

func TestLineContextHighlightWithMixedLanguages2(t *testing.T) {
	//    0         1         2         3         4         5         6
	//    0    5    0    5    0    5    0    5    0    5    0    5    0
	s := "the quick brown fox 我能吞下玻璃而不伤身体 jumps over the lazy dog"
	line := NewLineFromBytes(42, []byte(s))
	ctx := line.MarkLine(32, 37)

	got := ctx.Highlight("lorem ipsum")
	expected := strings.Join([]string{
		"  42:   the quick brown fox 我能吞下玻璃而不伤身体 jumps over the lazy dog",
		"                                                   ^^^^^",
		"                                                   lorem ipsum",
	}, "\n")
	if got != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, got)
	}
}
