package check

import (
	"slices"

	"github.com/flily/magi-c/ast"
)

type CheckList[T func(B) error, B any] struct {
	items []T
}

func NewCheckList[T func(B) error, B any](items ...T) *CheckList[T, B] {
	list := &CheckList[T, B]{
		items: slices.Clone(items),
	}

	return list
}

func (l *CheckList[T, B]) Check(node B) error {
	for _, check := range l.items {
		if err := check(node); err != nil {
			return err
		}
	}

	return nil
}

func checkDeclaration(d ast.Declaration) error {
	switch decl := d.(type) {
	case *ast.FunctionDeclaration:
		return checkFunctionDeclaration(decl)

	default:
		return d.Context().Error("unsupported declaration type %T", d)
	}
}

func checkDocument(doc *ast.Document) error {
	l := NewCheckList(
		checkDeclaration,
	)

	for _, decl := range doc.Declarations {
		if err := l.Check(decl); err != nil {
			return err
		}
	}

	return nil
}

type CodeChecker struct {
	document *ast.Document
}

func NewCodeChecker(document *ast.Document) *CodeChecker {
	c := &CodeChecker{
		document: document,
	}

	return c
}

func (c *CodeChecker) Check() error {
	l := NewCheckList(
		checkDocument,
	)

	return l.Check(c.document)
}
