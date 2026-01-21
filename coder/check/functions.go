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
			return ast.NewError(ectx, "duplicate function argument name: '%s'", name.Name)
		}

		nameMaps[name.Name] = name.Context()
	}

	return nil
}

func checkFunctionDeclaration(d *ast.FunctionDeclaration) error {
	l := NewCheckList(
		checkFunctionDeclarationNameDuplicate,
	)

	return l.Check(d)
}
