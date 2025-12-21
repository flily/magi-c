package csyntax

import (
	"strings"
)

type VariableDeclarator struct {
	PointerLevel int
	Name         string
	Initializer  Expression
}

type VariableDeclaration struct {
	Type       string
	Declarator []VariableDeclarator
}

func (v *VariableDeclaration) declarationNode() {}

func (v *VariableDeclaration) statementNode() {}

func (v *VariableDeclaration) Write(out *StyleWriter, level int) error {
	if err := out.WriteIndent(level); err != nil {
		return err
	}

	if err := out.Write("%s ", v.Type); err != nil {
		return err
	}

	for i, decl := range v.Declarator {
		if i > 0 {
			if err := out.Write(", "); err != nil {
				return err
			}
		}

		pointer := ""
		if decl.PointerLevel > 0 {
			pointer = strings.Repeat("*", decl.PointerLevel)
			if out.style.PointerSpacingBefore {
				pointer = " " + pointer
			}
			if out.style.PointerSpacingAfter {
				pointer = pointer + " "
			}
		}

		if err := out.Write("%s%s", pointer, decl.Name); err != nil {
			return err
		}

		if decl.Initializer != nil {
			if err := out.Write(" = "); err != nil {
				return err
			}
			if err := decl.Initializer.Write(out, level); err != nil {
				return err
			}
		}
	}

	return out.WriteLine(";")
}
