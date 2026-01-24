package context

import (
	"fmt"
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

func NewError(ctx *Context, message string, args ...any) *Diagnostic {
	m := fmt.Sprintf(message, args...)
	return Error.NewDiagnostic(ctx, m, "")
}

func (e *Diagnostic) Error() string {
	messageLine := fmt.Sprintf("%s: %s: %s", e.Context.PositionString(), e.Level, e.Message)
	return messageLine + DefaultNewLine + e.Context.HighlightText(e.Note)
}

func (e *Diagnostic) With(note string, args ...any) *Diagnostic {
	e.Note = fmt.Sprintf(note, args...)
	return e
}
