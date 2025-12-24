package csyntax

import (
	"fmt"
	"io"
	"strings"
)

type CStandard int

const (
	C89 CStandard = iota
	C99

	EOLCR   = "\r"
	EOLLF   = "\n"
	EOLCRLF = "\r\n"

	Space           = " "
	PointerAsterisk = "*"
	Semicolon       = ";"
	Comma           = ","
	Assign          = "="
	LeftBrace       = "{"
	RightBrace      = "}"
)

type CodeStyle struct {
	Indent                 string
	FunctionBraceOnNewLine bool
	FunctionBraceIndent    string
	IfBraceOnNewLine       bool
	IfBraceIndent          string
	ForBraceOnNewLine      bool
	ForBraceIndent         string
	WhileBraceOnNewLine    bool
	WhileBraceIndent       string
	SwitchBraceOnNewLine   bool
	SwitchBraceIndent      string
	CaseBranchIndent       string
	AssignmentSpacing      bool
	BinaryOperationSpacing bool
	TypeCastSpacing        bool
	CommaSpacingBefore     bool
	CommaSpacingAfter      bool
	PointerSpacingBefore   bool
	PointerSpacingAfter    bool
}

var (
	KRStyle = &CodeStyle{
		Indent:                 "    ",
		FunctionBraceOnNewLine: true,
		FunctionBraceIndent:    "",
		IfBraceOnNewLine:       false,
		IfBraceIndent:          "",
		ForBraceOnNewLine:      false,
		ForBraceIndent:         "",
		WhileBraceOnNewLine:    false,
		WhileBraceIndent:       "",
		SwitchBraceOnNewLine:   false,
		SwitchBraceIndent:      "",
		CaseBranchIndent:       "",
		AssignmentSpacing:      true,
		BinaryOperationSpacing: true,
		TypeCastSpacing:        true,
		CommaSpacingBefore:     false,
		CommaSpacingAfter:      true,
		PointerSpacingBefore:   false,
		PointerSpacingAfter:    true,
	}
)

func (s *CodeStyle) MakeWriter(out io.StringWriter) *StyleWriter {
	w := &StyleWriter{
		out:   out,
		style: s,
		EOL:   EOLLF,
	}

	return w
}

func (s *CodeStyle) Clone() *CodeStyle {
	clone := *s
	return &clone
}

func (s *CodeStyle) Comma() string {
	base := Comma

	if s.CommaSpacingBefore {
		base = Space + base
	}

	if s.CommaSpacingAfter {
		base = base + Space
	}

	return base
}

func (s *CodeStyle) Assign() string {
	if s.AssignmentSpacing {
		return Space + Assign + Space
	}

	return Assign
}

type Context struct {
	Filename string
	Line     int
}

func NewContext(filename string, line int) *Context {
	c := &Context{
		Filename: filename,
		Line:     line,
	}

	return c
}

func (c *Context) Write(out *StyleWriter, level int) error {
	return out.WriteLine("#line %d \"%s\"", c.Line, c.Filename)
}

type StyleWriter struct {
	out   io.StringWriter
	style *CodeStyle
	EOL   string
}

func (w *StyleWriter) Write(format string, args ...any) error {
	s := fmt.Sprintf(format, args...)
	_, err := w.out.WriteString(s)
	return err
}

func (w *StyleWriter) WriteLine(format string, args ...any) error {
	return w.Write(format+w.EOL, args...)
}

func (w *StyleWriter) WriteIndent(level int) error {
	indent := strings.Repeat(w.style.Indent, level)
	return w.Write("%s", indent)
}

type Node interface {
	Write(out *StyleWriter, level int) error
}

type Expression interface {
	Node
	expressionNode()
}

type Statement interface {
	Node
	statementNode()
}

type Declaration interface {
	Node
	declarationNode()
}
