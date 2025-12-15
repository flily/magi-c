package csyntax

import (
	"fmt"
	"io"
)

type CStandard int

const (
	C89 CStandard = iota
	C99

	EOLCR   = "\r"
	EOLLF   = "\n"
	EOLCRLF = "\r\n"
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

func (w *StyleWriter) WriteLine(format string, args ...any) error {
	s := fmt.Sprintf(format+w.EOL, args...)
	_, err := w.out.WriteString(s)
	return err
}

type Node interface {
	Write(out *StyleWriter, level int) error
}
