package csyntax

import (
	"strings"
)

type VariableDeclarator struct {
	PointerLevel int
	Name         string
	Initializer  Expression
}

func NewVariableDeclarator(name string, pointerLevel int, initializer Expression) VariableDeclarator {
	d := VariableDeclarator{
		PointerLevel: pointerLevel,
		Name:         name,
		Initializer:  initializer,
	}

	return d
}

type VariableDeclaration struct {
	Type       string
	Declarator []VariableDeclarator
}

func NewVariableDeclaration(typ string, declarators []VariableDeclarator) *VariableDeclaration {
	d := &VariableDeclaration{
		Type:       typ,
		Declarator: declarators,
	}

	return d
}

func (v *VariableDeclaration) declarationNode() {}

func (v *VariableDeclaration) statementNode() {}

// int a = 3;

func (v *VariableDeclaration) Write(out *StyleWriter, level int) error {
	if err := out.WriteIndent(level); err != nil {
		return err
	}

	if err := out.Write("%s ", v.Type); err != nil {
		return err
	}

	for i, decl := range v.Declarator {
		if i > 0 {
			if err := out.Write(out.style.Comma()); err != nil {
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
			if err := out.Write(out.style.Assign()); err != nil {
				return err
			}
			if err := decl.Initializer.Write(out, level); err != nil {
				return err
			}
		}
	}

	return out.WriteLine(";")
}

func (v *VariableDeclaration) Add(name string, pointerLevel int, initializer Expression) {
	decl := VariableDeclarator{
		Name:         name,
		PointerLevel: pointerLevel,
		Initializer:  initializer,
	}
	v.Declarator = append(v.Declarator, decl)
}
