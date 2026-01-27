package check

import (
	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
)

func checkFunctionDeclarationNameDuplicate(d *ast.FunctionDeclaration) error {
	nameMaps := make(map[string]*context.Context)
	for _, arg := range d.Arguments.Arguments {
		if arg.Name.IsDummy() {
			continue
		}

		name := arg.Name
		if ctx, found := nameMaps[name.Name]; found {
			ectx := name.Context()
			err := ectx.Error("duplicated function argument name: '%s'", name.Name).
				With("duplicated name").
				For(ctx.Note("first declared here"))
			return err
		}

		nameMaps[name.Name] = name.Context()
	}

	return nil
}

func checkFunctionReturnValue(d *ast.FunctionDeclaration) error {
	count := len(d.ReturnTypes.Types)

	retFound := false
	for _, stmt := range d.Statements {
		retStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			continue
		}

		retFound = true
		if len(retStmt.Value.Expressions) != count {
			c1 := d.ReturnTypes.Context()
			c2 := retStmt.Value.Context()
			err := c2.Error("function return value count mismatch, expect %d, got %d", count, len(retStmt.Value.Expressions)).
				With("SHALL return %d values", count).
				For(c1.Note("return value types is declared here"))
			return err
		}
	}

	if !retFound {
		c1 := d.ReturnTypes.Context()
		c2 := d.RBrace.Context()
		err := c2.Error("function missing return statement and reach the end of function").
			For(c1.Note("function return value types is declared here"))
		return err
	}

	return nil
}

func checkFunctionDeclaration(d *ast.FunctionDeclaration) error {
	l := NewCheckList(
		checkFunctionDeclarationNameDuplicate,
		checkFunctionReturnValue,
	)

	return l.Check(d)
}
