package tokenizer

import (
	"testing"

	"strings"

	"github.com/flily/magi-c/ast"
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

	word, err := tokenizer.ScanWord(0)
	if word == nil {
		t.Fatalf("expected a word token, got nil")
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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

func TestTokenizerScanSymbol(t *testing.T) {
	{
		buffer := []byte(strings.Join([]string{
			"====================",
		}, "\n"))

		tokenizer := NewTokenizerFrom(buffer, "test.txt")
		p1, err := tokenizer.ScanSymbol()
		exp1 := strings.Join([]string{
			"   1:   ====================",
			"        ^^^",
			"        here",
		}, "\n")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		got1 := p1.HighlightText("here")
		if got1 != exp1 {
			t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
		}
	}

	{
		buffer := []byte(strings.Join([]string{
			"0123456789",
		}, "\n"))

		tokenizer := NewTokenizerFrom(buffer, "test.txt")
		p1, err := tokenizer.ScanSymbol()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if p1 != nil {
			t.Errorf("expected nil, got %v", p1)
		}
	}
}

func TestTokenizerScanTokenOneSimpleLine(t *testing.T) {
	buffer := []byte(strings.Join([]string{
		"  a + b",
	}, "\n"))

	tokenizer := NewTokenizerFrom(buffer, "test.txt")

	ctxList := []ast.Node{}
	for i := range 3 {
		tokenizer.SkipWhitespace()
		tok, err := tokenizer.ScanToken()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if tok == nil {
			t.Fatalf("[%d] expected a token, got nil", i)
		}

		ctxList = append(ctxList, tok)
	}

	exp1 := strings.Join([]string{
		"   1:     a + b",
		"          ^",
		"          here",
	}, "\n")
	got1 := ctxList[0].HighlightText("here")
	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	exp2 := strings.Join([]string{
		"   1:     a + b",
		"            ^",
		"            here",
	}, "\n")
	got2 := ctxList[1].HighlightText("here")
	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	exp3 := strings.Join([]string{
		"   1:     a + b",
		"              ^",
		"              here",
	}, "\n")
	got3 := ctxList[2].HighlightText("here")
	if got3 != exp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp3, got3)
	}
}

func TestTokenizerScanTokenTwoSimpleLines(t *testing.T) {
	buffer := []byte(strings.Join([]string{
		"  aaaa + bbb",
		"ccc",
	}, "\n"))

	tokenizer := NewTokenizerFrom(buffer, "test.txt")

	ctxList := []ast.Node{}
	for i := range 4 {
		tokenizer.SkipWhitespace()
		tok, err := tokenizer.ScanToken()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if tok == nil {
			t.Fatalf("[%d] expected a token, got nil", i)
		}

		ctxList = append(ctxList, tok)
	}

	exp1 := strings.Join([]string{
		"   1:     aaaa + bbb",
		"          ^^^^",
		"          here",
	}, "\n")
	got1 := ctxList[0].HighlightText("here")
	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	if ctxList[0].Type() != ast.IdentifierName {
		t.Errorf("expected token type %s, got %s", ast.IdentifierName, ctxList[0].Type())
	}

	exp2 := strings.Join([]string{
		"   1:     aaaa + bbb",
		"               ^",
		"               here",
	}, "\n")
	got2 := ctxList[1].HighlightText("here")
	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	if ctxList[1].Type() != ast.Plus {
		t.Errorf("expected token type %s, got %s", ast.Plus, ctxList[1].Type())
	}

	exp3 := strings.Join([]string{
		"   1:     aaaa + bbb",
		"                 ^^^",
		"                 here",
	}, "\n")
	got3 := ctxList[2].HighlightText("here")
	if got3 != exp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp3, got3)
	}

	if ctxList[2].Type() != ast.IdentifierName {
		t.Errorf("expected token type %s, got %s", ast.IdentifierName, ctxList[2].Type())
	}

	exp4 := strings.Join([]string{
		"   2:   ccc",
		"        ^^^",
		"        here",
	}, "\n")
	got4 := ctxList[3].HighlightText("here")
	if got4 != exp4 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp4, got4)
	}

	if ctxList[3].Type() != ast.IdentifierName {
		t.Errorf("expected token type %s, got %s", ast.IdentifierName, ctxList[3].Type())
	}
}

func TestTokenizerScanTokenHexadecimalNumber(t *testing.T) {
	buffer := []byte(strings.Join([]string{
		"  0x1234 + 0xBEEF",
	}, "\n"))

	tokenizer := NewTokenizerFrom(buffer, "test.txt")

	ctxList := []ast.Node{}
	for i := range 3 {
		tokenizer.SkipWhitespace()
		tok, err := tokenizer.ScanToken()
		if err != nil {
			t.Fatalf("unexpected error:\n%v", err)
		}

		if tok == nil {
			t.Fatalf("[%d] expected a token, got nil", i)
		}

		ctxList = append(ctxList, tok)
	}

	exp1 := strings.Join([]string{
		"   1:     0x1234 + 0xBEEF",
		"          ^^^^^^",
		"          here",
	}, "\n")
	got1 := ctxList[0].HighlightText("here")
	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	if ctxList[0].Type() != ast.Integer {
		t.Errorf("expected token type %s, got %s", ast.Integer, ctxList[0].Type())
	}

	num1, ok := ctxList[0].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", ctxList[0])
	}

	if num1.Value != 0x1234 {
		t.Errorf("expected integer value 0x1234, got %d", num1.Value)
	}

	exp2 := strings.Join([]string{
		"   1:     0x1234 + 0xBEEF",
		"                 ^",
		"                 here",
	}, "\n")
	got2 := ctxList[1].HighlightText("here")
	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	if ctxList[1].Type() != ast.Plus {
		t.Errorf("expected token type %s, got %s", ast.Plus, ctxList[1].Type())
	}

	exp3 := strings.Join([]string{
		"   1:     0x1234 + 0xBEEF",
		"                   ^^^^^^",
		"                   here",
	}, "\n")
	got3 := ctxList[2].HighlightText("here")
	if got3 != exp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp3, got3)
	}

	if ctxList[2].Type() != ast.Integer {
		t.Errorf("expected token type %s, got %s", ast.Integer, ctxList[2].Type())
	}

	num3, ok := ctxList[2].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", ctxList[2])
	}

	if num3.Value != 0xbeef {
		t.Errorf("expected integer value 0xBEEF, got %d", num3.Value)
	}
}

func TestTokenizerScanTokenOctalNumber(t *testing.T) {
	buffer := []byte(strings.Join([]string{
		"  01234 + 0777",
	}, "\n"))

	tokenizer := NewTokenizerFrom(buffer, "test.txt")

	ctxList := []ast.Node{}
	for i := range 3 {
		tokenizer.SkipWhitespace()
		tok, err := tokenizer.ScanToken()
		if err != nil {
			t.Fatalf("unexpected error:\n%v", err)
		}

		if tok == nil {
			t.Fatalf("[%d] expected a token, got nil", i)
		}

		ctxList = append(ctxList, tok)
	}

	exp1 := strings.Join([]string{
		"   1:     01234 + 0777",
		"          ^^^^^",
		"          here",
	}, "\n")
	got1 := ctxList[0].HighlightText("here")
	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	if ctxList[0].Type() != ast.Integer {
		t.Errorf("expected token type %s, got %s", ast.Integer, ctxList[0].Type())
	}

	num1, ok := ctxList[0].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", ctxList[0])
	}

	if num1.Value != 01234 {
		t.Errorf("expected integer value 01234, got %d", num1.Value)
	}

	exp2 := strings.Join([]string{
		"   1:     01234 + 0777",
		"                ^",
		"                here",
	}, "\n")
	got2 := ctxList[1].HighlightText("here")
	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	if ctxList[1].Type() != ast.Plus {
		t.Errorf("expected token type %s, got %s", ast.Plus, ctxList[1].Type())
	}

	exp3 := strings.Join([]string{
		"   1:     01234 + 0777",
		"                  ^^^^",
		"                  here",
	}, "\n")
	got3 := ctxList[2].HighlightText("here")
	if got3 != exp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp3, got3)
	}

	if ctxList[2].Type() != ast.Integer {
		t.Errorf("expected token type %s, got %s", ast.Integer, ctxList[2].Type())
	}

	num3, ok := ctxList[2].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", ctxList[2])
	}

	if num3.Value != 0777 {
		t.Errorf("expected integer value 0777, got %d", num3.Value)
	}
}

func TestTokenizerScanDecimalInteger(t *testing.T) {
	buffer := []byte(strings.Join([]string{
		"  1234 + 5678",
	}, "\n"))

	tokenizer := NewTokenizerFrom(buffer, "test.txt")

	ctxList := []ast.Node{}
	for i := range 3 {
		tokenizer.SkipWhitespace()
		tok, err := tokenizer.ScanToken()
		if err != nil {
			t.Fatalf("unexpected error:\n%v", err)
		}

		if tok == nil {
			t.Fatalf("[%d] expected a token, got nil", i)
		}

		ctxList = append(ctxList, tok)
	}

	exp1 := strings.Join([]string{
		"   1:     1234 + 5678",
		"          ^^^^",
		"          here",
	}, "\n")
	got1 := ctxList[0].HighlightText("here")
	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	if ctxList[0].Type() != ast.Integer {
		t.Errorf("expected token type %s, got %s", ast.Integer, ctxList[0].Type())
	}

	num1, ok := ctxList[0].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", ctxList[0])
	}

	if num1.Value != 1234 {
		t.Errorf("expected integer value 1234, got %d", num1.Value)
	}

	exp2 := strings.Join([]string{
		"   1:     1234 + 5678",
		"               ^",
		"               here",
	}, "\n")
	got2 := ctxList[1].HighlightText("here")
	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	if ctxList[1].Type() != ast.Plus {
		t.Errorf("expected token type %s, got %s", ast.Plus, ctxList[1].Type())
	}

	exp3 := strings.Join([]string{
		"   1:     1234 + 5678",
		"                 ^^^^",
		"                 here",
	}, "\n")
	got3 := ctxList[2].HighlightText("here")
	if got3 != exp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp3, got3)
	}

	if ctxList[2].Type() != ast.Integer {
		t.Errorf("expected token Type %s, got %s", ast.Integer, ctxList[2].Type())
	}

	num3, ok := ctxList[2].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", ctxList[2])
	}

	if num3.Value != 5678 {
		t.Errorf("expected integer value 5678, got %d", num3.Value)
	}
}

func TestTokenizerScanDecimalFloat(t *testing.T) {
	buffer := []byte(strings.Join([]string{
		"  1234.5678 + 0.001 + 1.5e10 + 2.5E-3 + 3.5e+2",
	}, "\n"))

	tokenizer := NewTokenizerFrom(buffer, "test.txt")

	ctxList := []ast.Node{}
	for i := range 9 {
		tokenizer.SkipWhitespace()
		tok, err := tokenizer.ScanToken()
		if err != nil {
			t.Fatalf("unexpected error:\n%v", err)
		}

		if tok == nil {
			t.Fatalf("[%d] expected a token, got nil", i)
		}

		ctxList = append(ctxList, tok)
	}

	exp1 := strings.Join([]string{
		"   1:     1234.5678 + 0.001 + 1.5e10 + 2.5E-3 + 3.5e+2",
		"          ^^^^^^^^^",
		"          here",
	}, "\n")
	got1 := ctxList[0].HighlightText("here")
	if got1 != exp1 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
	}

	if ctxList[0].Type() != ast.Float {
		t.Errorf("expected token type %s, got %s", ast.Float, ctxList[0].Type())
	}

	num1, ok := ctxList[0].(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expected *ast.FloatLiteral, got %T", ctxList[0])
	}

	if num1.Value != 1234.5678 {
		t.Errorf("expected float value 1234.5678, got %f", num1.Value)
	}

	exp3 := strings.Join([]string{
		"   1:     1234.5678 + 0.001 + 1.5e10 + 2.5E-3 + 3.5e+2",
		"                      ^^^^^",
		"                      here",
	}, "\n")
	got3 := ctxList[2].HighlightText("here")
	if got3 != exp3 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp3, got3)
	}

	if ctxList[2].Type() != ast.Float {
		t.Errorf("expected token type %s, got %s", ast.Float, ctxList[2].Type())
	}

	num3, ok := ctxList[2].(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expected *ast.FloatLiteral, got %T", ctxList[2])
	}

	if num3.Value != 0.001 {
		t.Errorf("expected float value 0.001, got %f", num3.Value)
	}

	exp5 := strings.Join([]string{
		"   1:     1234.5678 + 0.001 + 1.5e10 + 2.5E-3 + 3.5e+2",
		"                              ^^^^^^",
		"                              here",
	}, "\n")
	got5 := ctxList[4].HighlightText("here")
	if got5 != exp5 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp5, got5)
	}

	if ctxList[4].Type() != ast.Float {
		t.Errorf("expected token type %s, got %s", ast.Float, ctxList[4].Type())
	}

	num5, ok := ctxList[4].(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expected *ast.FloatLiteral, got %T", ctxList[4])
	}

	if num5.Value != 1.5e10 {
		t.Errorf("expected float value 1.5e10, got %f", num5.Value)
	}

	exp7 := strings.Join([]string{
		"   1:     1234.5678 + 0.001 + 1.5e10 + 2.5E-3 + 3.5e+2",
		"                                       ^^^^^^",
		"                                       here",
	}, "\n")
	got7 := ctxList[6].HighlightText("here")
	if got7 != exp7 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp7, got7)
	}

	if ctxList[6].Type() != ast.Float {
		t.Errorf("expected token type %s, got %s", ast.Float, ctxList[6].Type())
	}

	num7, ok := ctxList[6].(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expected *ast.FloatLiteral, got %T", ctxList[6])
	}

	if num7.Value != 2.5e-3 {
		t.Errorf("expected float value 2.5e-3, got %f", num7.Value)
	}

	exp9 := strings.Join([]string{
		"   1:     1234.5678 + 0.001 + 1.5e10 + 2.5E-3 + 3.5e+2",
		"                                                ^^^^^^",
		"                                                here",
	}, "\n")
	got9 := ctxList[8].HighlightText("here")
	if got9 != exp9 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp9, got9)
	}

	if ctxList[8].Type() != ast.Float {
		t.Errorf("expected token type %s, got %s", ast.Float, ctxList[8].Type())
	}

	num9, ok := ctxList[8].(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expected *ast.FloatLiteral, got %T", ctxList[8])
	}

	if num9.Value != 3.5e2 {
		t.Errorf("expected float value 3.5e2, got %f", num9.Value)
	}
}
