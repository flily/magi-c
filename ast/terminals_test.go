package ast

import (
	"testing"

	"strings"
)

func TestStringLiteral(t *testing.T) {
	text := "lorem ipsum"
	ctxList := generateTestWords(text)

	s := NewStringLiteral(ctxList[0], "lorem")
	if s.Type() != String {
		t.Fatalf("string literal type expected %d, got %d", String, s.Type())
	}

	var _ TerminalNode = s
	var _ Expression = s
	s.expressionNode()

	a := ASTBuildValue("lorem")

	if err := s.EqualTo(nil, a); err != nil {
		t.Fatalf("expected string literal not equal to actual:\n%s", err)
	}
}

func TestStringLiteralNotEqual(t *testing.T) {
	text := "lorem ipsum"
	ctxList := generateTestWords(text)

	s := NewStringLiteral(ctxList[0], "lorem")
	{
		a := ASTBuildValue(3)
		exp := strings.Join([]string{
			"   1:   lorem ipsum",
			"        ^^^^^",
			"        expect a *ast.IntegerLiteral",
		}, "\n")

		err := s.EqualTo(s, a)
		if err == nil {
			t.Fatalf("expected string literal not equal to actual")
		}

		if err.Error() != exp {
			t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", exp, err.Error())
		}
	}

	{
		a := ASTBuildValue("ipsum")
		err := s.EqualTo(s, a)
		if err == nil {
			t.Fatalf("expected string literal not equal to actual")
		}

		exp := strings.Join([]string{
			"   1:   lorem ipsum",
			"        ^^^^^",
			"        wrong string value, expect 'ipsum', got 'lorem'",
		}, "\n")
		if err.Error() != exp {
			t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", exp, err.Error())
		}
	}
}

func TestIntegerLiteral(t *testing.T) {
	text := "1234 5678"
	ctxList := generateTestWords(text)

	i := NewIntegerLiteral(ctxList[0], 1234)
	if i.Type() != Integer {
		t.Fatalf("integer literal type expected %d, got %d", Integer, i.Type())
	}

	var _ TerminalNode = i
	var _ Expression = i
	i.expressionNode()

	a := ASTBuildValue(1234)

	if err := i.EqualTo(nil, a); err != nil {
		t.Fatalf("expected integer literal not equal to actual:\n%s", err)
	}
}

func TestIntegerLiteralUnsigned(t *testing.T) {
	text := "1234 5678"
	ctxList := generateTestWords(text)

	u := NewIntegerLiteral(ctxList[0], 1234)
	if u.Type() != Integer {
		t.Fatalf("integer literal type expected %d, got %d", Integer, u.Type())
	}

	var _ TerminalNode = u
	var _ Expression = u
	u.expressionNode()

	a := ASTBuildValue(uint64(1234))

	if err := u.EqualTo(nil, a); err != nil {
		t.Fatalf("expected integer literal not equal to actual:\n%s", err)
	}
}

func TestIntegerLiteralNotEqual(t *testing.T) {
	text := "1234 5678"
	ctxList := generateTestWords(text)

	i := NewIntegerLiteral(ctxList[0], 1234)
	{
		a := ASTBuildValue("1234")
		exp := strings.Join([]string{
			"   1:   1234 5678",
			"        ^^^^",
			"        expect a *ast.StringLiteral",
		}, "\n")

		err := i.EqualTo(i, a)
		if err == nil {
			t.Fatalf("expected integer literal not equal to actual")
		}

		if err.Error() != exp {
			t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", exp, err.Error())
		}
	}

	{
		a := ASTBuildValue(5678)
		err := i.EqualTo(i, a)
		if err == nil {
			t.Fatalf("expected integer literal not equal to actual")
		}

		exp := strings.Join([]string{
			"   1:   1234 5678",
			"        ^^^^",
			"        wrong integer value, expect 5678, got 1234",
		}, "\n")
		if err.Error() != exp {
			t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", exp, err.Error())
		}
	}
}

func TestFloatLiteral(t *testing.T) {
	text := "3.14 2.71"
	ctxList := generateTestWords(text)

	f := NewFloatLiteral(ctxList[0], 3.14)
	if f.Type() != Float {
		t.Fatalf("float literal type expected %d, got %d", Float, f.Type())
	}

	var _ TerminalNode = f
	var _ Expression = f
	f.expressionNode()

	a := ASTBuildValue(3.14)

	if err := f.EqualTo(nil, a); err != nil {
		t.Fatalf("expected float literal not equal to actual:\n%s", err)
	}
}

func TestFloatLiteralNotEqual(t *testing.T) {
	text := "3.14 2.71"
	ctxList := generateTestWords(text)

	f := NewFloatLiteral(ctxList[0], 3.14)
	{
		a := ASTBuildValue("3.14")
		exp := strings.Join([]string{
			"   1:   3.14 2.71",
			"        ^^^^",
			"        expect a *ast.StringLiteral",
		}, "\n")

		err := f.EqualTo(f, a)
		if err == nil {
			t.Fatalf("expected float literal not equal to actual")
		}

		if err.Error() != exp {
			t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", exp, err.Error())
		}
	}

	{
		a := ASTBuildValue(2.71)
		err := f.EqualTo(f, a)
		if err == nil {
			t.Fatalf("expected float literal not equal to actual")
		}

		exp := strings.Join([]string{
			"   1:   3.14 2.71",
			"        ^^^^",
			"        wrong float value, expect 2.71, got 3.14",
		}, "\n")
		if err.Error() != exp {
			t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", exp, err.Error())
		}
	}
}

func TestBuildValueUnsupported(t *testing.T) {
	exp := "ASTBuildValue: unsupported value type: *int"

	errMsg := ""
	defer func() {
		if r := recover(); r != nil {
			errMsg = r.(string)
		}

		if errMsg != exp {
			t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", exp, errMsg)
		}
	}()

	a := 0
	_ = ASTBuildValue(&a)
}

func TestIdentifier(t *testing.T) {
	text := "lorem ipsum"
	ctxList := generateTestWords(text)

	id := NewIdentifier(ctxList[0])
	if id.Type() != IdentifierName {
		t.Fatalf("identifier type expected %d, got %d", IdentifierName, id.Type())
	}

	var _ TerminalNode = id
	var _ Expression = id
	id.expressionNode()

	a := ASTBuildIdentifier("lorem")

	if err := id.EqualTo(nil, a); err != nil {
		t.Fatalf("expected identifier not equal to actual:\n%s", err)
	}
}

func TestIdentifierNotEqual(t *testing.T) {
	text := "lorem ipsum"
	ctxList := generateTestWords(text)

	id := NewIdentifier(ctxList[0])
	{
		a := ASTBuildValue(1234)
		exp := strings.Join([]string{
			"   1:   lorem ipsum",
			"        ^^^^^",
			"        expect a *ast.IntegerLiteral",
		}, "\n")

		err := id.EqualTo(id, a)
		if err == nil {
			t.Fatalf("expected identifier not equal to actual")
		}

		if err.Error() != exp {
			t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", exp, err.Error())
		}
	}

	{
		a := ASTBuildIdentifier("ipsum")
		err := id.EqualTo(id, a)
		if err == nil {
			t.Fatalf("expected identifier not equal to actual")
		}

		exp := strings.Join([]string{
			"   1:   lorem ipsum",
			"        ^^^^^",
			"        wrong identifier name, expect ipsum, got lorem",
		}, "\n")
		if err.Error() != exp {
			t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", exp, err.Error())
		}
	}
}

func TestTerminalToken(t *testing.T) {
	text := "+ if"
	ctxList := generateTestWords(text)

	symbol := NewTerminalToken(ctxList[0], Plus)
	if symbol.Type() != Plus {
		t.Fatalf("terminal token type expected %d, got %d", Plus, symbol.Type())
	}

	var _ TerminalNode = symbol

	a := ASTBuildSymbol(Plus)

	if err := symbol.EqualTo(nil, a); err != nil {
		t.Fatalf("expected terminal token not equal to actual:\n%s", err)
	}

	keyword := NewTerminalToken(ctxList[1], If)
	if keyword.Type() != If {
		t.Fatalf("terminal token type expected %d, got %d", If, keyword.Type())
	}

	var _ TerminalNode = keyword

	b := ASTBuildKeyword(If)

	if err := keyword.EqualTo(nil, b); err != nil {
		t.Fatalf("expected terminal token not equal to actual:\n%s", err)
	}
}

func TestTerminalTokenNotEqual(t *testing.T) {
	text := "+ if"
	ctxList := generateTestWords(text)

	symbol := NewTerminalToken(ctxList[0], Plus)
	{
		a := ASTBuildSymbol(Sub)
		exp := strings.Join([]string{
			"   1:   + if",
			"        ^",
			"        wrong token type, expect '-', got '+'",
		}, "\n")

		err := symbol.EqualTo(symbol, a)
		if err == nil {
			t.Fatalf("expected terminal token not equal to actual")
		}

		if err.Error() != exp {
			t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", exp, err.Error())
		}
	}

	keyword := NewTerminalToken(ctxList[1], If)
	{
		b := ASTBuildKeyword(Else)
		exp := strings.Join([]string{
			"   1:   + if",
			"          ^^",
			"          wrong token type, expect 'else', got 'if'",
		}, "\n")

		err := keyword.EqualTo(keyword, b)
		if err == nil {
			t.Fatalf("expected terminal token not equal to actual")
		}

		if err.Error() != exp {
			t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", exp, err.Error())
		}
	}

	{
		c := ASTBuildValue(42)
		exp := strings.Join([]string{
			"   1:   + if",
			"        ^",
			"        expect a *ast.IntegerLiteral",
		}, "\n")

		err := symbol.EqualTo(symbol, c)
		if err == nil {
			t.Fatalf("expected terminal token not equal to actual")
		}

		if err.Error() != exp {
			t.Fatalf("wrong error message:\nexpected:\n%s\ngot:\n%s", exp, err.Error())
		}
	}
}
