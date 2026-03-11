package csyntax

// Base on ISO/IEC 9899:1999

import (
	"fmt"
	"io"
	"strings"

	"github.com/flily/magi-c/context"
)

type CodeStyle struct {
	Indent                 StringElement
	FunctionBraceOnNewLine StyleBoolean
	FunctionBraceIndent    StringElement
	IfSpacing              StyleBoolean
	IfBraceOnNewLine       StyleBoolean
	IfBraceIndent          StringElement
	ForBraceOnNewLine      StyleBoolean
	ForBraceIndent         StringElement
	WhileSpacing           StyleBoolean
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
		FunctionBraceOnNewLine: true,
		FunctionBraceIndent:    "",
		IfSpacing:              true,
		IfBraceOnNewLine:       false,
		IfBraceIndent:          "",
		ForBraceOnNewLine:      false,
		ForBraceIndent:         "",
		WhileSpacing:           true,
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

func (s *CodeStyle) GetIndent(level Level) StringElement {
	return StringElement(strings.Repeat(string(s.Indent), level.IndentLevel))
}

func (s *CodeStyle) IfNewLine(level Level) ElementCollection {
	result := []CodeElement{
		DelimiterSpace,
	}

	if s.IfBraceOnNewLine {
		result = []CodeElement{
			s.EOL,
			s.IfBraceIndent,
			s.GetIndent(level),
		}
	}

	return result
}

func (s *CodeStyle) WhileNewLine(level Level) ElementCollection {
	result := []CodeElement{
		DelimiterSpace,
	}

	if s.WhileBraceOnNewLine {
		result = []CodeElement{
			s.EOL,
			s.WhileBraceIndent,
			s.GetIndent(level),
		}
	}

	return result
}

func (s *CodeStyle) BinaryOperator(op Punctuator) ElementCollection {
	result := []CodeElement{
		s.BinaryOperationSpacing.Select(DelimiterSpace),
		op,
		s.BinaryOperationSpacing.Select(DelimiterSpace),
	}

	return result
}

type Context struct {
	Context *context.Context
}

func NewContext(ctx *context.Context) *Context {
	c := &Context{
		Context: ctx,
	}

	return c
}

func (c *Context) codeElement()     {}
func (c *Context) declarationNode() {}
func (c *Context) statementNode()   {}

func (c *Context) Write(out *StyleWriter, level Level) error {
	filename, line, _ := c.Context.Position()
	return out.Write(level,
		PreprocessorLine, DelimiterSpace, NewIntegerStringElement(line+1),
		DelimiterSpace, StringElement("\""+filename+"\""), out.style.EOL,
	)
}

type CodeElement interface {
	codeElement()
}

type Node interface {
	CodeElement
	Write(out *StyleWriter, level Level) error
}

type Statement interface {
	Node
	statementNode()
}

type Declaration interface {
	Node
	declarationNode()
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

func (c ElementCollection) Write(out *StyleWriter, level Level) error {
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
	Write(out *StyleWriter, level Level) error
}

type StyleWriter struct {
	out              io.StringWriter
	style            *CodeStyle
	lastWasDelimiter bool
}

func (w *StyleWriter) WriteIndent(level Level) error {
	return w.Write(level, w.style.GetIndent(level))
}

func (w *StyleWriter) Write(level Level, items ...CodeElement) error {
	for i, item := range items {
		_, isDelimiter := item.(DelimiterCharacter)
		switch it := item.(type) {
		case StringElement:
			if _, err := w.out.WriteString(it.ItemString()); err != nil {
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

func (w *StyleWriter) WriteLine(level Level, items ...CodeElement) error {
	err := w.Write(level, items...)
	if err != nil {
		return err
	}

	_, err = w.out.WriteString(w.style.EOL.String())
	return err
}

func (w *StyleWriter) WriteIndentLine(level Level, items ...CodeElement) error {
	if err := w.WriteIndent(level); err != nil {
		return err
	}

	return w.WriteLine(level, items...)
}
