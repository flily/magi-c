package tokenizer

import (
	"testing"

	"strings"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
	"github.com/flily/magi-c/preprocessor"
)

func checkContext(t *testing.T, ctx *context.Context, expected string) {
	t.Helper()

	got := ctx.HighlightText("here")
	if got != expected {
		t.Errorf("context got wrong output, expect:\n%s\ngot:\n%s", expected, got)
	}
}

func checkTerminalNode(t *testing.T, node ast.Node, expectedType ast.TokenType, expectedText string) {
	t.Helper()

	if node == nil {
		t.Fatalf("expect a non-nil token")
	}

	term, ok := node.(ast.TerminalNode)
	if !ok {
		t.Fatalf("expected TerminalNode, got %T", node)
	}

	if term.Type() != expectedType {
		t.Errorf("expect token type %s, got %s", expectedType, term.Type())
	}

	got := term.HighlightText("here")
	if got != expectedText {
		t.Errorf("expect:\n%s\ngot:\n%s", expectedText, got)
	}
}

func checkError(t *testing.T, err error, expected string) {
	t.Helper()

	got := err.Error()
	if got != expected {
		t.Errorf("expect error:\n%s\ngot:\n%s", expected, got)
	}
}

func TestTokenizerScanAll(t *testing.T) {
	code := strings.Join([]string{
		"#include <stdio.h>",
		"",
		"fun main() {",
		"}",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")
	tokenizer.RegisterPreprocessor("include", preprocessor.Include)
	tokens, err := tokenizer.ScanAll()
	if err != nil {
		t.Fatalf("unexpected error:\n%v", err)
	}

	if len(tokens) != 7 {
		t.Fatalf("expected 7 tokens, got %d", len(tokens))
	}

	directiveTypes := []ast.TokenType{
		ast.NodePreprocessorInclude,
		ast.Function,
		ast.IdentifierName,
		ast.LeftParen,
		ast.RightParen,
		ast.LeftBrace,
		ast.RightBrace,
	}

	for i, expectedType := range directiveTypes {
		term := tokens[i]
		if term.Type() != expectedType {
			t.Errorf("token %d: expected type %s, got %s", i, expectedType, term.Type())
		}
	}
}

func TestTokenizerScanEOF(t *testing.T) {
	code := strings.Join([]string{
		"#include <stdio.h>",
		"",
		"fun main() {",
		"}",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")
	tokenizer.RegisterPreprocessor("include", preprocessor.Include)
	tokens, err := tokenizer.ScanAll()
	if err != nil {
		t.Fatalf("unexpected error:\n%v", err)
	}

	if len(tokens) != 7 {
		t.Fatalf("expected 7 tokens, got %d", len(tokens))
	}

	last, err := tokenizer.ScanToken()
	if err != nil {
		t.Fatalf("unexpected error when scanning after EOF:\n%v", err)
	}

	if last != nil {
		t.Fatalf("expected last token to be nil, got %v", last)
	}

	exp := strings.Join([]string{
		"   4:   }<EOF>",
		"         ^^^^^",
		"         here",
	}, "\n")
	checkContext(t, tokenizer.EOFContext(), exp)
}

func TestTokenizerSkipWhitespace(t *testing.T) {
	code := strings.Join([]string{
		"                lorem ipsum",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

	_, p1 := tokenizer.CurrentChar()
	exp1 := strings.Join([]string{
		"   1:                   lorem ipsum",
		"        ^",
		"        here",
	}, "\n")
	checkContext(t, p1, exp1)

	tokenizer.SkipWhitespace()
	_, p2 := tokenizer.CurrentChar()
	exp2 := strings.Join([]string{
		"   1:                   lorem ipsum",
		"                        ^",
		"                        here",
	}, "\n")
	checkContext(t, p2, exp2)
}

func TestTokenizerSkipWhitespaceToNextLine(t *testing.T) {
	code := strings.Join([]string{
		"        lorem        ",
		"    ",
		"        \t\t\t        ",
		"      ipsum dolor sit amet",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

	_, p1 := tokenizer.CurrentChar()
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

	_, p2 := tokenizer.CurrentChar()
	exp2 := strings.Join([]string{
		"   1:           lorem        ",
		"                ^",
		"                here",
	}, "\n")

	got2 := p2.HighlightText("here")
	if got2 != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got2)
	}

	word := tokenizer.ScanWordToken(0)
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

	_, p3 := tokenizer.CurrentChar()
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

	_, p4 := tokenizer.CurrentChar()
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
	code := strings.Join([]string{
		"====================",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

	_, p1 := tokenizer.CurrentChar()
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
	_, p2 := tokenizer.CurrentChar()
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
	_, p3 := tokenizer.CurrentChar()
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
		code := strings.Join([]string{
			"====================",
		}, "\n")

		tokenizer := NewTokenizerFromString(code, "test.txt")
		p1, err := tokenizer.ScanSymbol()
		exp1 := strings.Join([]string{
			"   1:   ====================",
			"        ^^^",
			"        here",
		}, "\n")

		if err != nil {
			t.Fatalf("unexpected error:\n%v", err)
		}

		got1 := p1.HighlightText("here")
		if got1 != exp1 {
			t.Errorf("expected:\n%s\ngot:\n%s", exp1, got1)
		}
	}

	{
		code := strings.Join([]string{
			"0123456789",
		}, "\n")

		tokenizer := NewTokenizerFromString(code, "test.txt")
		exp := strings.Join([]string{
			"   1:   0123456789",
			"        ^",
			"        invalid symbol '0'",
		}, "\n")

		p1, err := tokenizer.ScanSymbol()
		if p1 != nil {
			t.Fatalf("expected nil, got %v", p1)
		}

		got := err.Error()
		if got != exp {
			t.Errorf("expected:\n%s\ngot:\n%s", exp, got)
		}
	}
}

func TestTokenizerScanTokenOneSimpleLine(t *testing.T) {
	code := strings.Join([]string{
		"  a + b",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

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
	code := strings.Join([]string{
		"  aaaa + bbb",
		"ccc",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

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
	checkTerminalNode(t, ctxList[0], ast.IdentifierName, exp1)

	exp2 := strings.Join([]string{
		"   1:     aaaa + bbb",
		"               ^",
		"               here",
	}, "\n")
	checkTerminalNode(t, ctxList[1], ast.Plus, exp2)

	exp3 := strings.Join([]string{
		"   1:     aaaa + bbb",
		"                 ^^^",
		"                 here",
	}, "\n")
	checkTerminalNode(t, ctxList[2], ast.IdentifierName, exp3)

	exp4 := strings.Join([]string{
		"   2:   ccc",
		"        ^^^",
		"        here",
	}, "\n")
	checkTerminalNode(t, ctxList[3], ast.IdentifierName, exp4)
}

func TestTokenizerScanTokenHexadecimalNumber(t *testing.T) {
	code := strings.Join([]string{
		"  0x1234 + 0xBEEF + 0xc0de",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

	ctxList := []ast.Node{}
	for i := range 5 {
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
		"   1:     0x1234 + 0xBEEF + 0xc0de",
		"          ^^^^^^",
		"          here",
	}, "\n")
	checkTerminalNode(t, ctxList[0], ast.Integer, exp1)

	num1, ok := ctxList[0].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", ctxList[0])
	}

	if num1.Value != 0x1234 {
		t.Errorf("expected integer value 0x1234, got %d", num1.Value)
	}

	exp2 := strings.Join([]string{
		"   1:     0x1234 + 0xBEEF + 0xc0de",
		"                 ^",
		"                 here",
	}, "\n")
	checkTerminalNode(t, ctxList[1], ast.Plus, exp2)

	exp3 := strings.Join([]string{
		"   1:     0x1234 + 0xBEEF + 0xc0de",
		"                   ^^^^^^",
		"                   here",
	}, "\n")
	checkTerminalNode(t, ctxList[2], ast.Integer, exp3)

	num3, ok := ctxList[2].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", ctxList[2])
	}

	if num3.Value != 0xbeef {
		t.Errorf("expected integer value 0xBEEF, got %d", num3.Value)
	}

	exp4 := strings.Join([]string{
		"   1:     0x1234 + 0xBEEF + 0xc0de",
		"                          ^",
		"                          here",
	}, "\n")
	checkTerminalNode(t, ctxList[3], ast.Plus, exp4)

	exp5 := strings.Join([]string{
		"   1:     0x1234 + 0xBEEF + 0xc0de",
		"                            ^^^^^^",
		"                            here",
	}, "\n")
	checkTerminalNode(t, ctxList[4], ast.Integer, exp5)

	num5, ok := ctxList[4].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", ctxList[4])
	}

	if num5.Value != 0xc0de {
		t.Errorf("expected integer value 0xc0de, got %d", num5.Value)
	}
}

func TestTokenizerScanTokenHexadecimalNumberErrorNoNumberEOL(t *testing.T) {
	code := strings.Join([]string{
		"  0x",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

	tokenizer.SkipWhitespace()
	result, err := tokenizer.ScanToken()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}

	exp := strings.Join([]string{
		"   1:     0x",
		"          ^^",
		"          invalid hexadecimal number '0x'",
	}, "\n")
	checkError(t, err, exp)
}

func TestTokenizerScanTokenHexadecimalNumberErrorInvalidFormat(t *testing.T) {
	code := strings.Join([]string{
		"  0xGHIJ+0xghij",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")
	tokenizer.SkipWhitespace()
	result, err := tokenizer.ScanToken()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}

	exp := strings.Join([]string{
		"   1:     0xGHIJ+0xghij",
		"          ^^^^^^",
		"          invalid hexadecimal number '0xGHIJ'",
	}, "\n")
	checkError(t, err, exp)
}

func TestTokenizerScanTokenHexadecimalNumberErrorTooLargeNumber(t *testing.T) {
	code := strings.Join([]string{
		"  0x01234567890ABCDEF1234 + 0xbeef",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")
	tokenizer.SkipWhitespace()
	result, err := tokenizer.ScanToken()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}

	exp := strings.Join([]string{
		"   1:     0x01234567890ABCDEF1234 + 0xbeef",
		"          ^^^^^^^^^^^^^^^^^^^^^^^",
		"          hexadecimal number '0x01234567890ABCDEF1234' is too large",
	}, "\n")
	checkError(t, err, exp)
}

func TestTokenizerScanTokenOctalNumber(t *testing.T) {
	code := strings.Join([]string{
		"  01234 + 0777",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

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
	checkTerminalNode(t, ctxList[0], ast.Integer, exp1)

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
	checkTerminalNode(t, ctxList[1], ast.Plus, exp2)

	exp3 := strings.Join([]string{
		"   1:     01234 + 0777",
		"                  ^^^^",
		"                  here",
	}, "\n")
	checkTerminalNode(t, ctxList[2], ast.Integer, exp3)

	num3, ok := ctxList[2].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", ctxList[2])
	}

	if num3.Value != 0777 {
		t.Errorf("expected integer value 0777, got %d", num3.Value)
	}
}

func TestTokenizerScanTokenOctalNumberErrorInvalidFormat(t *testing.T) {
	code := strings.Join([]string{
		"  0123456789",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

	tokenizer.SkipWhitespace()
	result, err := tokenizer.ScanToken()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}

	exp := strings.Join([]string{
		"   1:     0123456789",
		"          ^^^^^^^^^^",
		"          invalid octal number '0123456789'",
	}, "\n")
	checkError(t, err, exp)
}

func TestTokenizerScanTokenOctalNumberErrorTooLargeNumber(t *testing.T) {
	code := strings.Join([]string{
		"  01234567012345670123456701234567",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

	tokenizer.SkipWhitespace()
	result, err := tokenizer.ScanToken()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}

	exp := strings.Join([]string{
		"   1:     01234567012345670123456701234567",
		"          ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^",
		"          octal number '01234567012345670123456701234567' is too large",
	}, "\n")
	checkError(t, err, exp)
}

func TestTokenizerScanTokenDecimalInteger(t *testing.T) {
	code := strings.Join([]string{
		"  1234 + 5678",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

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
	checkTerminalNode(t, ctxList[0], ast.Integer, exp1)

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
	checkTerminalNode(t, ctxList[1], ast.Plus, exp2)

	exp3 := strings.Join([]string{
		"   1:     1234 + 5678",
		"                 ^^^^",
		"                 here",
	}, "\n")
	checkTerminalNode(t, ctxList[2], ast.Integer, exp3)

	num3, ok := ctxList[2].(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", ctxList[2])
	}

	if num3.Value != 5678 {
		t.Errorf("expected integer value 5678, got %d", num3.Value)
	}
}

func TestTokenizerScanTokenDecimalSingleZero(t *testing.T) {
	code := strings.Join([]string{
		"  0",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")
	tokenizer.SkipWhitespace()
	tok, err := tokenizer.ScanToken()
	if err != nil {
		t.Fatalf("unexpected error:\n%v", err)
	}

	exp1 := strings.Join([]string{
		"   1:     0",
		"          ^",
		"          here",
	}, "\n")
	checkTerminalNode(t, tok, ast.Integer, exp1)

	num, ok := tok.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expected *ast.IntegerLiteral, got %T", tok)
	}

	if num.Value != 0 {
		t.Errorf("expected integer value 0, got %d", num.Value)
	}

	_, afterPos := tokenizer.CurrentChar()
	exp2 := strings.Join([]string{
		"   1:     0<EOF>",
		"           ^^^^^",
		"           here",
	}, "\n")
	got := afterPos.HighlightText("here")
	if got != exp2 {
		t.Errorf("expected:\n%s\ngot:\n%s", exp2, got)
	}
}

func TestTokenizerScanTokenDecimalNumberFloat(t *testing.T) {
	code := strings.Join([]string{
		"  1234.5678 + 0.001 + 1.5e10 + 2.5E-3 + 3e+2",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

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
		"   1:     1234.5678 + 0.001 + 1.5e10 + 2.5E-3 + 3e+2",
		"          ^^^^^^^^^",
		"          here",
	}, "\n")
	checkTerminalNode(t, ctxList[0], ast.Float, exp1)

	num1, ok := ctxList[0].(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expected *ast.FloatLiteral, got %T", ctxList[0])
	}

	if num1.Value != 1234.5678 {
		t.Errorf("expected float value 1234.5678, got %f", num1.Value)
	}

	exp3 := strings.Join([]string{
		"   1:     1234.5678 + 0.001 + 1.5e10 + 2.5E-3 + 3e+2",
		"                      ^^^^^",
		"                      here",
	}, "\n")
	checkTerminalNode(t, ctxList[2], ast.Float, exp3)

	num3, ok := ctxList[2].(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expected *ast.FloatLiteral, got %T", ctxList[2])
	}

	if num3.Value != 0.001 {
		t.Errorf("expected float value 0.001, got %f", num3.Value)
	}

	exp5 := strings.Join([]string{
		"   1:     1234.5678 + 0.001 + 1.5e10 + 2.5E-3 + 3e+2",
		"                              ^^^^^^",
		"                              here",
	}, "\n")
	checkTerminalNode(t, ctxList[4], ast.Float, exp5)

	num5, ok := ctxList[4].(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expected *ast.FloatLiteral, got %T", ctxList[4])
	}

	if num5.Value != 1.5e10 {
		t.Errorf("expected float value 1.5e10, got %f", num5.Value)
	}

	exp7 := strings.Join([]string{
		"   1:     1234.5678 + 0.001 + 1.5e10 + 2.5E-3 + 3e+2",
		"                                       ^^^^^^",
		"                                       here",
	}, "\n")
	checkTerminalNode(t, ctxList[6], ast.Float, exp7)

	num7, ok := ctxList[6].(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expected *ast.FloatLiteral, got %T", ctxList[6])
	}

	if num7.Value != 2.5e-3 {
		t.Errorf("expected float value 2.5e-3, got %f", num7.Value)
	}

	exp9 := strings.Join([]string{
		"   1:     1234.5678 + 0.001 + 1.5e10 + 2.5E-3 + 3e+2",
		"                                                ^^^^",
		"                                                here",
	}, "\n")
	checkTerminalNode(t, ctxList[8], ast.Float, exp9)

	num9, ok := ctxList[8].(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("expected *ast.FloatLiteral, got %T", ctxList[8])
	}

	if num9.Value != 3e2 {
		t.Errorf("expected float value 3e2, got %f", num9.Value)
	}
}

func TestTokenizerScanTokenDecimalNumberFloatErrorInvalidFormatInIntegerPart(t *testing.T) {
	code := strings.Join([]string{
		"  123dfg + 3.14",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

	tokenizer.SkipWhitespace()
	result, err := tokenizer.ScanToken()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}

	exp := strings.Join([]string{
		"   1:     123dfg + 3.14",
		"          ^^^^^^",
		"          invalid decimal number '123dfg'",
	}, "\n")
	checkError(t, err, exp)
}

func TestTokenizerScanTokenDecimalNumberFloatErrorInvalidFormatInFractionPart(t *testing.T) {
	code := strings.Join([]string{
		"  3.14xyz + 123",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

	tokenizer.SkipWhitespace()
	result, err := tokenizer.ScanToken()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}

	exp := strings.Join([]string{
		"   1:     3.14xyz + 123",
		"          ^^^^^^^",
		"          invalid decimal number '3.14xyz'",
	}, "\n")
	checkError(t, err, exp)
}

func TestTokenizerScanTokenDecimalNumberFloatErrorInvalidFormatInExponentPart(t *testing.T) {
	code := strings.Join([]string{
		"  1.5e10xyz + 123",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")

	tokenizer.SkipWhitespace()
	result, err := tokenizer.ScanToken()
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}

	exp := strings.Join([]string{
		"   1:     1.5e10xyz + 123",
		"          ^^^^^^^^^",
		"          invalid decimal number '1.5e10xyz'",
	}, "\n")
	checkError(t, err, exp)
}

func TestTokenizerScanTokenPreprocessorDirective(t *testing.T) {
	code := strings.Join([]string{
		"  #include <stdio.h>",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")
	tokenizer.RegisterPreprocessor("include", preprocessor.Include)

	tokenizer.SkipWhitespace()
	tok, err := tokenizer.ScanToken()
	if err != nil {
		t.Fatalf("unexpected error:\n%v", err)
	}

	if tok == nil {
		t.Fatalf("expected a token, got nil")
	}

	exp := strings.Join([]string{
		"   1:     #include <stdio.h>",
		"          ^^^^^^^^ ^^^^^^^^^",
		"          here",
	}, "\n")
	checkTerminalNode(t, tok, ast.NodePreprocessorInclude, exp)
}

func TestTokenizerScanSimplestProgram(t *testing.T) {
	code := strings.Join([]string{
		"fun main() {",
		"  return 0",
		"}",
	}, "\n")

	tokenizer := NewTokenizerFromString(code, "test.txt")
	nodes, err := tokenizer.ScanAll()
	if err != nil {
		t.Fatalf("unexpected error:\n%v", err)
	}

	if len(nodes) != 8 {
		t.Fatalf("expected 8 nodes, got %d", len(nodes))
	}
}
