package token

import (
	"testing"

	"strings"
)

func TestTokenizerSkipWhitespace(t *testing.T) {
	buffer := []byte(strings.Join([]string{
		"                lorem ipsum",
	}, "\n"))

	tokenizer := NewTokenizerFrom(buffer, "test.txt")

	p1 := tokenizer.CurrentChar()
	exp1 := strings.Join([]string{
		"   1:                   lorem ipsum",
		"        ^",
		"        here",
	}, "\n")

	got1 := p1.HighlightText("here")
	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	tokenizer.SkipWhitespace()

	p2 := tokenizer.CurrentChar()
	exp2 := strings.Join([]string{
		"   1:                   lorem ipsum",
		"                        ^",
		"                        here",
	}, "\n")

	got2 := p2.HighlightText("here")
	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}
}

func TestTokenizerSkipWhitespaceToNextLine(t *testing.T) {
	buffer := []byte(strings.Join([]string{
		"        lorem        ",
		"    ",
		"        \t\t\t        ",
		"      ipsum dolor sit amet",
	}, "\n"))

	tokenizer := NewTokenizerFrom(buffer, "test.txt")

	p1 := tokenizer.CurrentChar()
	exp1 := strings.Join([]string{
		"   1:           lorem        ",
		"        ^",
		"        here",
	}, "\n")

	got1 := p1.HighlightText("here")
	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	tokenizer.SkipWhitespace()

	p2 := tokenizer.CurrentChar()
	exp2 := strings.Join([]string{
		"   1:           lorem        ",
		"                ^",
		"                here",
	}, "\n")

	got2 := p2.HighlightText("here")
	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	word := tokenizer.ScanWord()
	if word == nil {
		t.Fatalf("expected a word token, got nil")
	}

	expWord := strings.Join([]string{
		"   1:           lorem        ",
		"                ^^^^^",
		"                here",
	}, "\n")

	gotWord := word.HighlightText("here")
	if gotWord != expWord {
		t.Errorf("expected:\n%s\ngot:\n%s", expWord, gotWord)
	}

	p3 := tokenizer.CurrentChar()
	exp3 := strings.Join([]string{
		"   1:           lorem        ",
		"                     ^",
		"                     here",
	}, "\n")

	got3 := p3.HighlightText("here")
	if got3 != exp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp3, got3)
	}

	tokenizer.SkipWhitespace()

	p4 := tokenizer.CurrentChar()
	exp4 := strings.Join([]string{
		"   4:         ipsum dolor sit amet",
		"              ^",
		"              here",
	}, "\n")

	got4 := p4.HighlightText("here")
	if got4 != exp4 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp4, got4)
	}
}

func TestTokenizerScanFixedString(t *testing.T) {
	buffer := []byte(strings.Join([]string{
		"====================",
	}, "\n"))

	tokenizer := NewTokenizerFrom(buffer, "test.txt")

	p1 := tokenizer.CurrentChar()
	exp1 := strings.Join([]string{
		"   1:   ====================",
		"        ^",
		"        here",
	}, "\n")

	got1 := p1.HighlightText("here")
	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	t2 := tokenizer.ScanFixedString("=")
	p2 := tokenizer.CurrentChar()
	expt2 := strings.Join([]string{
		"   1:   ====================",
		"        ^",
		"        here",
	}, "\n")

	gott2 := t2.HighlightText("here")
	if gott2 != expt2 {
		t.Errorf("expected:\n%s\ngot:\n%s", expt2, gott2)
	}

	expp2 := strings.Join([]string{
		"   1:   ====================",
		"         ^",
		"         here",
	}, "\n")
	gotp2 := p2.HighlightText("here")
	if gotp2 != expp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", expp2, gotp2)
	}

	t3 := tokenizer.ScanFixedString("==")
	p3 := tokenizer.CurrentChar()
	expt3 := strings.Join([]string{
		"   1:   ====================",
		"         ^^",
		"         here",
	}, "\n")

	gott3 := t3.HighlightText("here")
	if gott3 != expt3 {
		t.Errorf("expected:\n%s\ngot:\n%s", expt3, gott3)
	}

	expp3 := strings.Join([]string{
		"   1:   ====================",
		"           ^",
		"           here",
	}, "\n")
	gotp3 := p3.HighlightText("here")
	if gotp3 != expp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", expp3, gotp3)
	}
}
