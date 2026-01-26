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
			ectx := context.Join(ctx, name.Context())
			return ectx.Error("duplicate function argument name: '%s'", name.Name).With("duplicate function argument name")
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
			c2 := retStmt.Context()
			ectx := context.Join(c1, c2)
			return ectx.Error("function return value count mismatch, expect %d, got %d", count, len(retStmt.Value.Expressions))
		}
	}

	if !retFound {
		c1 := d.ReturnTypes.Context()
		c2 := d.RBrace.Context()
		ectx := context.Join(c1, c2)
		return ectx.Error("function missing return statement")
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
