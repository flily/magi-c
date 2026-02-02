package check

import (
	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/context"
)

func checkFunctionDeclarationNameDuplicate(d *ast.FunctionDeclaration) context.DiagnosticInfo {
	if d.Arguments == nil {
		return nil
	}

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

func checkFunctionReturnValue(d *ast.FunctionDeclaration) context.DiagnosticInfo {
	if d.ReturnTypes == nil {
		return nil
	}

	count := d.ReturnTypes.Length()

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

func checkFunctionMainDeclaration(d *ast.FunctionDeclaration) context.DiagnosticInfo {
	if d.Name.Name != "main" {
		return nil
	}

	if d.ReturnTypes == nil || d.ReturnTypes.Length() == 0 {
		// fun main() { ... }
		return nil
	}

	if d.ReturnTypes.Length() != 1 {
		ctx := d.ReturnTypes.Context()
		err := ctx.Error("function 'main' must have return type 'int' or no return type, got %d return types", d.ReturnTypes.Length()).
			With("int or no return type")
		return err
	}

	// TODO: return type check
	return nil
}

func checkFunctionDeclaration(conf *CheckConfigure, d *ast.FunctionDeclaration) *context.DiagnosticContainer {
	l := NewCheckList(
		checkFunctionDeclarationNameDuplicate,
		checkFunctionReturnValue,
		checkFunctionMainDeclaration,
	)

	return l.Check(conf, d)
}
