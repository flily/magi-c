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
)

type StyleBoolean bool

func (b StyleBoolean) Select(value CodeElement) CodeElement {
	if b {
		return value
	}

	return DelimiterNone
}

type CodeStyle struct {
	Indent                 StringElement
	FunctionBraceOnNewLine StyleBoolean
	FunctionBraceIndent    StringElement
	IfBraceOnNewLine       StyleBoolean
	IfBraceIndent          StringElement
	ForBraceOnNewLine      StyleBoolean
	ForBraceIndent         StringElement
	WhileBraceOnNewLine    StyleBoolean
	WhileBraceIndent       StringElement
	SwitchBraceOnNewLine   StyleBoolean
	SwitchBraceIndent      StringElement
	CaseBranchIndent       StringElement
	AssignmentSpacing      StyleBoolean
	BinaryOperationSpacing StyleBoolean
	TypeCastSpacing        StyleBoolean
	CommaSpacingBefore     StyleBoolean
	CommaSpacingAfter      StyleBoolean
	PointerSpacingBefore   StyleBoolean
	PointerSpacingAfter    StyleBoolean
	EOL                    StringElement
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
		EOL:                    EOLLF,
	}
)

func (s *CodeStyle) MakeWriter(out io.StringWriter) *StyleWriter {
	w := &StyleWriter{
		out:   out,
		style: s,
	}

	return w
}

func (s *CodeStyle) Clone() *CodeStyle {
	clone := *s
	return &clone
}

func (s *CodeStyle) Comma() ElementCollection {
	result := []CodeElement{
		s.CommaSpacingBefore.Select(DelimiterSpace),
		PunctuatorComma,
		s.CommaSpacingAfter.Select(DelimiterSpace),
	}

	return result
}

func (s *CodeStyle) Assign() ElementCollection {
	result := []CodeElement{
		s.AssignmentSpacing.Select(DelimiterSpace),
		OperatorAssign,
		s.AssignmentSpacing.Select(DelimiterSpace),
	}

	return result
}

func (s *CodeStyle) FunctionNewLine() ElementCollection {
	result := []CodeElement{
		DelimiterSpace,
	}

	if s.FunctionBraceOnNewLine {
		result = []CodeElement{
			s.EOL,
			s.FunctionBraceIndent,
		}
	}

	return result
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

func (c *Context) codeElement() {}

func (c *Context) Write(out *StyleWriter, level int) error {
	return out.Write(level,
		PreprocessorLine, DelimiterSpace, NewIntegerStringElement(c.Line),
		DelimiterSpace, StringElement("\""+c.Filename+"\""), out.style.EOL,
	)
}

type CodeElement interface {
	codeElement()
}

type ElementCollection []CodeElement

func (c ElementCollection) codeElement() {}

func FromCodeElements[T CodeElement](items ...T) ElementCollection {
	collection := make(ElementCollection, len(items))
	for i, item := range items {
		collection[i] = item
	}

	return collection
}

func NewElementCollection(items ...CodeElement) ElementCollection {
	return FromCodeElements(items...)
}

func (c ElementCollection) Write(out *StyleWriter, level int) error {
	return out.Write(level, c...)
}

func (c ElementCollection) On(cond bool) ElementCollection {
	if cond {
		return c
	}

	return nil
}

func (c ElementCollection) Select(cond bool, alt ElementCollection) ElementCollection {
	if cond {
		return c
	}

	return alt
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

func NewIntegerStringElement(i int) StringElement {
	return StringElement(fmt.Sprintf("%d", i))
}

func FormatStringElement(format string, args ...any) StringElement {
	s := fmt.Sprintf(format, args...)
	return StringElement(s)
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
	DelimiterNone    StringElement      = ""
	DelimiterSpace   DelimiterCharacter = " "
	DefaultDelimiter DelimiterCharacter = " "
)

type StyleWriter struct {
	out              io.StringWriter
	style            *CodeStyle
	lastWasDelimiter bool
}

func (w *StyleWriter) MakeIndent(level int) StringElement {
	return StringElement(strings.Repeat(string(w.style.Indent), level))
}

func (w *StyleWriter) WriteIndent(level int) error {
	return w.Write(0, w.MakeIndent(level))
}

func (w *StyleWriter) Write(level int, items ...CodeElement) error {
	for i, item := range items {
		_, isDelimiter := item.(DelimiterCharacter)
		switch it := item.(type) {
		case StringElement:
			if _, err := w.out.WriteString(it.String()); err != nil {
				return err
			}
			w.lastWasDelimiter = isDelimiter

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
	}

	return nil
}

func (w *StyleWriter) WriteLine(level int, items ...CodeElement) error {
	err := w.Write(level, items...)
	if err != nil {
		return err
	}

	_, err = w.out.WriteString(w.style.EOL.String())
	return err
}

func (w *StyleWriter) WriteIndentLine(level int, items ...CodeElement) error {
	if err := w.WriteIndent(level); err != nil {
		return err
	}

	return w.WriteLine(0, items...)
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
