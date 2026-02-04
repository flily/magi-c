package check

import (
	"slices"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
)

type CheckConfigure struct {
	Level context.ErrorLevel
}

func NewDefaultCheckConfigure() *CheckConfigure {
	c := &CheckConfigure{
		Level: context.Error,
	}

	return c
}

type CheckFunc[B any] func(*CheckConfigure, B) context.DiagnosticInfo

type CheckList[T func(B) context.DiagnosticInfo, B any] struct {
	items []T
}

func NewCheckList[T func(B) context.DiagnosticInfo, B any](items ...T) *CheckList[T, B] {
	list := &CheckList[T, B]{
		items: slices.Clone(items),
	}

	return list
}

func (l *CheckList[T, B]) Check(conf *CheckConfigure, node B) *context.DiagnosticContainer {
	c := context.NewDiagnosticContainer(conf.Level)

	for _, check := range l.items {
		err := check(node)
		if err != nil {
			e := c.Add(err)
			if e != nil {
				return c
			}
		}
	}

	return c
}

type CheckRunner[T func(*CheckConfigure, B) *context.DiagnosticContainer, B any] struct {
	items []T
}

func NewCheckRunner[T func(*CheckConfigure, B) *context.DiagnosticContainer, B any](items ...T) *CheckRunner[T, B] {
	runner := &CheckRunner[T, B]{
		items: slices.Clone(items),
	}

	return runner
}

func (r *CheckRunner[T, B]) Run(conf *CheckConfigure, node B) *context.DiagnosticContainer {
	c := context.NewDiagnosticContainer(conf.Level)

	for _, check := range r.items {
		err := check(conf, node)
		if err != nil {
			e := c.Merge(err)
			if e != nil {
				return c
			}
		}
	}

	return c
}

func checkDeclaration(conf *CheckConfigure, d ast.Declaration) *context.DiagnosticContainer {
	switch decl := d.(type) {
	case *ast.FunctionDeclaration:
		return checkFunctionDeclaration(conf, decl)

	case *ast.PreprocessorInclude:
		return nil

	case *ast.PreprocessorInline:
		return nil

	default:
		return d.Context().Error("unsupported declaration type %T", d).ToContainer()
	}
}

func checkDocument(conf *CheckConfigure, doc *ast.Document) *context.DiagnosticContainer {
	l := NewCheckRunner(
		checkDeclaration,
	)

	c := context.NewDiagnosticContainer(conf.Level)
	for _, decl := range doc.Declarations {
		err := l.Run(conf, decl)
		if err != nil {
			e := c.Merge(err)
			if e != nil {
				return c
			}
		}
	}

	return c
}

type CodeChecker struct {
	config   *CheckConfigure
	document *ast.Document
}

func NewCodeChecker(conf *CheckConfigure, document *ast.Document) *CodeChecker {
	c := &CodeChecker{
		config:   conf,
		document: document,
	}

	return c
}

func (c *CodeChecker) Check() *context.DiagnosticContainer {
	l := NewCheckRunner(
		checkDocument,
	)

	return l.Run(c.config, c.document)
}
