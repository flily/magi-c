package context

import (
	"fmt"
	"strings"
)

type ErrorLevel int

const (
	Ignored ErrorLevel = iota
	Note
	Remark
	Warning
	Error
	Fatal
)

var levelNames = map[ErrorLevel]string{
	Ignored: "ignored",
	Note:    "note",
	Remark:  "remark",
	Warning: "warning",
	Error:   "error",
	Fatal:   "fatal",
}

func (l ErrorLevel) String() string {
	return levelNames[l]
}

func (l ErrorLevel) NewDiagnostic(ctx *Context, message string, note string) *Diagnostic {
	return NewDiagnostic(l, ctx, message, note)
}

type Diagnostic struct {
	Level   ErrorLevel
	Message string
	Context *Context
	Note    string
}

func NewDiagnostic(level ErrorLevel, ctx *Context, message string, note string) *Diagnostic {
	e := &Diagnostic{
		Level:   level,
		Message: message,
		Context: ctx,
		Note:    note,
	}

	return e
}

func NewNote(ctx *Context, message string, args ...any) *Diagnostic {
	m := fmt.Sprintf(message, args...)
	return Note.NewDiagnostic(ctx, m, "")
}

func NewWarning(ctx *Context, message string, args ...any) *Diagnostic {
	m := fmt.Sprintf(message, args...)
	return Warning.NewDiagnostic(ctx, m, "")
}

func NewError(ctx *Context, message string, args ...any) *Diagnostic {
	m := fmt.Sprintf(message, args...)
	return Error.NewDiagnostic(ctx, m, "")
}

func (d *Diagnostic) Error() string {
	messageLine := fmt.Sprintf("%s: %s: %s", d.Context.PositionString(), d.Level, d.Message)
	return messageLine + DefaultNewLine + d.Context.HighlightText(d.Note)
}

func (d *Diagnostic) With(note string, args ...any) *Diagnostic {
	d.Note = fmt.Sprintf(note, args...)
	return d
}

func (d *Diagnostic) For(reason *Diagnostic) *DiagnosticCombo {
	c := &DiagnosticCombo{
		Info:   d,
		Reason: reason,
	}

	return c
}

type DiagnosticCombo struct {
	Info   *Diagnostic
	Reason *Diagnostic
}

func (c *DiagnosticCombo) Error() string {
	parts := make([]string, 0, 2)
	parts = append(parts, c.Info.Error())

	if c.Reason != nil {
		parts = append(parts, c.Reason.Error())
	}

	return strings.Join(parts, DefaultNewLine)
}
