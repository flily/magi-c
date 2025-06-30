package token

import (
	"testing"
)

func TestContextHighlightWithoutLineNo(t *testing.T) {
	text := "the quick brown fox jumps over the lazy dog"
	line := NewLineContextFromString("example.txt", -1, text)
	line.AddHighlight(10, 15)
	ctx := NewContext("example.txt", line)

	if line.Content.Valid() {
		t.Errorf("LineContext should not be valid when line number is -1")
	}

	expected := JoinLines([]string{
		"the quick brown fox jumps over the lazy dog",
		"          ^^^^^",
		"          lorem ipsum",
	})

	got := ctx.Highlight("lorem ipsum")
	if got != expected {
		t.Errorf("wrong highlighted result, got:\n%s\nexpect:\n%s", got, expected)
	}
}

func TestContextHighlightChineseWithoutLineNo(t *testing.T) {
	text := "我能吞下玻璃而不伤身体"
	line := NewLineContextFromString("example.txt", -1, text)
	line.AddHighlight(4, 6)
	ctx := NewContext("example.txt", line)

	expected := JoinLines([]string{
		"我能吞下玻璃而不伤身体",
		"        ^^^^",
		"        lorem ipsum",
	})

	got := ctx.Highlight("lorem ipsum")
	if got != expected {
		t.Errorf("wrong highlighted result, got:\n%s\nexpect:\n%s", got, expected)
	}
}

func TestContextHighlightWithLineNo(t *testing.T) {
	text := "the quick brown fox jumps over the lazy dog"
	line := NewLineContextFromString("example.txt", 41, text)
	line.AddHighlight(10, 15)
	ctx := NewContext("example.txt", line)

	expected := JoinLines([]string{
		"42    the quick brown fox jumps over the lazy dog",
		"                ^^^^^",
		"                lorem ipsum",
	})

	got := ctx.Highlight("lorem ipsum")
	if got != expected {
		t.Errorf("wrong highlighted result, got:\n%s\nexpect:\n%s", got, expected)
	}
}

func TestContextHighlightMultipleParts(t *testing.T) {
	text := "the quick brown fox jumps over the lazy dog"
	line := NewLineContextFromString("example.txt", 41, text)
	line.AddHighlight(4, 9)
	line.AddHighlight(10, 15)
	line.AddHighlight(35, 39)
	ctx := NewContext("example.txt", line)

	expected := JoinLines([]string{
		"42    the quick brown fox jumps over the lazy dog",
		"          ^^^^^ ^^^^^                    ^^^^",
		"          lorem ipsum",
	})

	got := ctx.Highlight("lorem ipsum")
	if got != expected {
		t.Errorf("wrong highlighted result, got:\n%s\nexpect:\n%s", got, expected)
	}
}

func TestContextHighlightPreviousLinesWithLineNo(t *testing.T) {
	text := "the quick brown fox jumps over the lazy dog"
	previous := []*LineContext{
		NewLineContextFromString("example.txt", 37, "11111111"),
		NewLineContextFromString("example.txt", 38, "22222222"),
		NewLineContextFromString("example.txt", 39, "33333333"),
		NewLineContextFromString("example.txt", 40, "44444444"),
	}

	line := NewLineContextFromString("example.txt", 41, text)
	line.AddHighlight(10, 15)
	ctx := NewContext("example.txt", line)
	ctx.AddPrevLines(previous)

	expected := JoinLines([]string{
		"38    11111111",
		"39    22222222",
		"40    33333333",
		"41    44444444",
		"42    the quick brown fox jumps over the lazy dog",
		"                ^^^^^",
		"                lorem ipsum",
	})

	got := ctx.Highlight("lorem ipsum")
	if got != expected {
		t.Errorf("wrong highlighted result, got:\n%s\nexpect:\n%s", got, expected)
	}
}

func TestContextHighlightMultipleLinesWithLineNo(t *testing.T) {
	text1 := "the quick brown fox"
	text2 := "jumps over the lazy dog"
	previous := []*LineContext{
		NewLineContextFromString("example.txt", 37, "11111111"),
		NewLineContextFromString("example.txt", 38, "22222222"),
		NewLineContextFromString("example.txt", 39, "33333333"),
		NewLineContextFromString("example.txt", 40, "44444444"),
	}

	texts := []*LineContext{
		NewLineContextFromString("example.txt", 41, text1),
		NewLineContextFromString("example.txt", 42, text2),
	}

	texts[0].AddHighlight(10, 15)
	texts[1].AddHighlight(15, 19)

	ctx1 := NewContext("example.txt", texts...)
	ctx1.AddPrevLines(previous)

	expected := JoinLines([]string{
		"38    11111111",
		"39    22222222",
		"40    33333333",
		"41    44444444",
		"42    the quick brown fox",
		"                ^^^^^",
		"43    jumps over the lazy dog",
		"                     ^^^^",
		"                     lorem ipsum",
	})

	got := ctx1.Highlight("lorem ipsum")
	if got != expected {
		t.Errorf("wrong highlighted result, got:\n%s\nexpect:\n%s", got, expected)
	}
}

func TestContextHighlightWithTabIndent(t *testing.T) {
	text := "  \tthe  quick \tbrown fox jumps over the lazy dog"
	previous := []*LineContext{
		NewLineContextFromString("example.txt", 37, "11111111"),
		NewLineContextFromString("example.txt", 38, "22222222"),
		NewLineContextFromString("example.txt", 39, "33333333"),
		NewLineContextFromString("example.txt", 40, "44444444"),
	}

	line := NewLineContextFromString("example.txt", 41, text)
	line.AddHighlight(15, 20)
	ctx := NewContext("example.txt", line)
	ctx.AddPrevLines(previous)

	expected := JoinLines([]string{
		"38    11111111",
		"39    22222222",
		"40    33333333",
		"41    44444444",
		"42      	the  quick 	brown fox jumps over the lazy dog",
		"        	           	^^^^^",
		"        	           	lorem ipsum",
	})

	got := ctx.Highlight("lorem ipsum")
	if got != expected {
		t.Errorf("wrong highlighted result, got:\n%s\nexpect:\n%s", got, expected)
	}
}

func TestJoinMultipleLineContexts(t *testing.T) {

	previous1 := []*LineContext{
		NewLineContextFromString("example.txt", 37, "11111111"),
		NewLineContextFromString("example.txt", 38, "22222222"),
		NewLineContextFromString("example.txt", 39, "33333333"),
		NewLineContextFromString("example.txt", 40, "44444444"),
	}

	line1 := NewLineContextFromString("example.txt", 41, "the quick brown fox")
	line1.AddHighlight(10, 15)
	ctx1 := NewContext("example.txt", line1)
	ctx1.AddPrevLines(previous1)

	previous2 := []*LineContext{
		NewLineContextFromString("example.txt", 37, "11111111"),
		NewLineContextFromString("example.txt", 38, "22222222"),
		NewLineContextFromString("example.txt", 39, "33333333"),
		NewLineContextFromString("example.txt", 40, "44444444"),
		NewLineContextFromString("example.txt", 41, "the quick brown fox"),
	}

	line2 := NewLineContextFromString("example.txt", 42, "jumps over the lazy dog")
	line2.AddHighlight(15, 19)
	ctx2 := NewContext("example.txt", line2)
	ctx2.AddPrevLines(previous2)

	result := JoinContexts(ctx1, ctx2)
	if result == nil {
		t.Fatalf("JoinContexts returned nil")
	}

	expected := JoinLines([]string{
		"38    11111111",
		"39    22222222",
		"40    33333333",
		"41    44444444",
		"42    the quick brown fox",
		"                ^^^^^",
		"43    jumps over the lazy dog",
		"                     ^^^^",
		"                     lorem ipsum",
	})

	got := result.Highlight("lorem ipsum")
	if got != expected {
		t.Errorf("wrong highlighted result, got:\n%s\nexpect:\n%s", got, expected)
	}
}

func TestJoinMultipleLineContextsWithHighlightsInSameLine(t *testing.T) {

	previous1 := []*LineContext{
		NewLineContextFromString("example.txt", 37, "11111111"),
		NewLineContextFromString("example.txt", 38, "22222222"),
		NewLineContextFromString("example.txt", 39, "33333333"),
		NewLineContextFromString("example.txt", 40, "44444444"),
	}

	line1 := NewLineContextFromString("example.txt", 41, "the quick brown fox")
	line1.AddHighlight(10, 15)
	ctx1 := NewContext("example.txt", line1)
	ctx1.AddPrevLines(previous1)

	previous2 := []*LineContext{
		NewLineContextFromString("example.txt", 37, "11111111"),
		NewLineContextFromString("example.txt", 38, "22222222"),
		NewLineContextFromString("example.txt", 39, "33333333"),
		NewLineContextFromString("example.txt", 40, "44444444"),
	}

	line2 := NewLineContextFromString("example.txt", 41, "the quick brown fox")
	line2.AddHighlight(4, 9)
	line3 := NewLineContextFromString("example.txt", 42, "jumps over the lazy dog")
	line3.AddHighlight(15, 19)
	ctx2 := NewContext("example.txt", line2, line3)
	ctx2.AddPrevLines(previous2)

	result := JoinContexts(ctx1, ctx2)
	if result == nil {
		t.Fatalf("JoinContexts returned nil")
	}

	expected := JoinLines([]string{
		"38    11111111",
		"39    22222222",
		"40    33333333",
		"41    44444444",
		"42    the quick brown fox",
		"          ^^^^^ ^^^^^",
		"43    jumps over the lazy dog",
		"                     ^^^^",
		"                     lorem ipsum",
	})

	got := result.Highlight("lorem ipsum")
	if got != expected {
		t.Errorf("wrong highlighted result, got:\n%s\nexpect:\n%s", got, expected)
	}
}
