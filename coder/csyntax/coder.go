package csyntax

// Base on ISO/IEC 9899:1999

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
	Indent                 StringElement
	FunctionBraceOnNewLine bool
	FunctionBraceIndent    StringElement
	IfBraceOnNewLine       bool
	IfBraceIndent          StringElement
	ForBraceOnNewLine      bool
	ForBraceIndent         StringElement
	WhileBraceOnNewLine    bool
	WhileBraceIndent       StringElement
	SwitchBraceOnNewLine   bool
	SwitchBraceIndent      StringElement
	CaseBranchIndent       StringElement
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
		FunctionBraceOnNewLine: false,
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
	return out.WriteItems(level, PreprocessorLine, DelimiterSpace, NewIntegerElement(c.Line),
		DelimiterSpace, NewStringElement("\""+c.Filename+"\""), out.EOL)
}

type CodeElement interface {
	codeElement()
}

type ElementCollection []CodeElement

func (c ElementCollection) codeElement() {}

func (c ElementCollection) Write(out *StyleWriter, level int) error {
	return out.WriteItems(level, c...)
}

func NewElementCollection(items ...CodeElement) ElementCollection {
	return ElementCollection(items)
}

type WritableItem interface {
	CodeElement
	ItemString() string
}

type WritableCollection interface {
	CodeElement
	Write(out *StyleWriter, level int) error
}

type StringElement string

func NewStringElement(s string) StringElement {
	return StringElement(s)
}

func NewIntegerElement(i int) StringElement {
	return StringElement(fmt.Sprintf("%d", i))
}

func (e StringElement) codeElement() {}

func (e StringElement) String() string {
	return string(e)
}

func (e StringElement) ItemString() string {
	return string(e)
}

type DelimiterCharacter string

func NewDelimiter(c string) DelimiterCharacter {
	return DelimiterCharacter(c)
}

func (d DelimiterCharacter) codeElement() {}

func (d DelimiterCharacter) String() string {
	return string(d)
}

func (d DelimiterCharacter) ItemString() string {
	return string(d)
}

const (
	DelimiterNone                       = ""
	DelimiterSpace   DelimiterCharacter = " "
	DefaultDelimiter DelimiterCharacter = " "
)

type StyleWriter struct {
	out              io.StringWriter
	style            *CodeStyle
	EOL              StringElement
	lastWasDelimiter bool
}

func (w *StyleWriter) Write(format string, args ...any) error {
	s := fmt.Sprintf(format, args...)
	_, err := w.out.WriteString(s)
	return err
}

func (w *StyleWriter) WriteLine(format string, args ...any) error {
	return w.Write(format+string(w.EOL), args...)
}

func (w *StyleWriter) MakeIndent(level int) StringElement {
	return StringElement(strings.Repeat(string(w.style.Indent), level))
}

func (w *StyleWriter) WriteIndent(level int) error {
	return w.WriteItems(0, w.MakeIndent(level))
}

func (w *StyleWriter) WriteItems(level int, items ...CodeElement) error {
	for i, item := range items {
		_, isDelimiter := item.(DelimiterCharacter)

		switch it := item.(type) {
		case StringElement:
			if _, err := w.out.WriteString(it.String()); err != nil {
				return err
			}

		case WritableItem:
			if w.lastWasDelimiter && isDelimiter {
				continue
			}

			w.lastWasDelimiter = isDelimiter
			if _, err := w.out.WriteString(it.ItemString()); err != nil {
				return err
			}

		case WritableCollection:
			if err := it.Write(w, level); err != nil {
				return err
			}

		default:
			err := fmt.Sprintf("type %T in %d is not acceptable", item, i)
			panic(err)
		}

		w.lastWasDelimiter = isDelimiter
	}

	return nil
}

type Node interface {
	CodeElement
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
