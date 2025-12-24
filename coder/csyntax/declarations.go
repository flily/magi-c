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

func (v *VariableDeclaration) Write(out *StyleWriter, level int) error {
	if err := out.WriteIndent(level); err != nil {
		return err
	}

	if err := out.Write("%s", v.Type); err != nil {
		return err
	}

	for i, decl := range v.Declarator {
		if i > 0 {
			if err := out.Write(out.style.Comma()); err != nil {
				return err
			}
		}

		commaSpace := out.style.CommaSpacingAfter

		pointer := ""
		if decl.PointerLevel > 0 {
			pointer = strings.Repeat(PointerAsterisk, decl.PointerLevel)
			if !commaSpace && (i == 0 && out.style.PointerSpacingBefore) {
				pointer = Space + pointer
			}
			if out.style.PointerSpacingAfter {
				pointer = pointer + Space
			}
		}

		if len(pointer) <= 0 && i == 0 {
			// space between type and first declarator
			pointer = Space
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

	return out.WriteLine(Semicolon)
}

func (v *VariableDeclaration) Add(name string, pointerLevel int, initializer Expression) {
	decl := NewVariableDeclarator(name, pointerLevel, initializer)
	v.Declarator = append(v.Declarator, decl)
}
