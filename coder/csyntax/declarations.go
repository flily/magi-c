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

type ParameterListItem struct {
	Type *Type
	Name string
}

func NewParameterListItem(typ *Type, name string) *ParameterListItem {
	p := &ParameterListItem{
		Type: typ,
		Name: name,
	}

	return p
}

func (i *ParameterListItem) Write(out *StyleWriter, level int) error {
	if err := i.Type.Write(out, level); err != nil {
		return err
	}

	if i.Type.PointerLevel <= 0 || (!out.style.PointerSpacingAfter && i.Type.PointerLevel > 0) {
		if err := out.Write(Space); err != nil {
			return err
		}
	}

	return out.Write("%s", i.Name)
}

type ParameterList struct {
	Items []*ParameterListItem
}

func NewParameterList(items ...*ParameterListItem) *ParameterList {
	p := &ParameterList{
		Items: items,
	}

	return p
}

func (p *ParameterList) Write(out *StyleWriter, level int) error {
	for i, item := range p.Items {
		if i > 0 {
			if err := out.Write(out.style.Comma()); err != nil {
				return err
			}
		}

		if err := item.Write(out, level); err != nil {
			return err
		}
	}

	return nil
}

type FunctionDeclaration struct {
	ReturnType *Type
	Name       string
	Parameters *ParameterList
	Body       []Statement
}

func NewFunctionDeclaration(name string, returnType *Type, parameters *ParameterList, body []Statement) *FunctionDeclaration {
	f := &FunctionDeclaration{
		ReturnType: returnType,
		Name:       name,
		Parameters: parameters,
		Body:       body,
	}

	return f
}

func (f *FunctionDeclaration) declarationNode() {}

func (f *FunctionDeclaration) AddStatement(stmt Statement) {
	f.Body = append(f.Body, stmt)
}

func (f *FunctionDeclaration) Write(out *StyleWriter, level int) error {
	if err := out.WriteIndent(level); err != nil {
		return err
	}

	if err := f.ReturnType.Write(out, level); err != nil {
		return err
	}

	if err := out.Write(" %s(", f.Name); err != nil {
		return err
	}

	if err := f.Parameters.Write(out, level); err != nil {
		return err
	}

	if err := out.Write(")"); err != nil {
		return err
	}

	if out.style.FunctionBraceOnNewLine {
		if err := out.WriteLine("%s%s%s", out.EOL, out.style.FunctionBraceIndent, LeftBrace); err != nil {
			return err
		}
	} else {
		if err := out.WriteLine(" %s", LeftBrace); err != nil {
			return err
		}
	}

	for _, stmt := range f.Body {
		if err := stmt.Write(out, level+1); err != nil {
			return err
		}
	}

	if err := out.WriteLine("%s%s", out.style.FunctionBraceIndent, RightBrace); err != nil {
		return err
	}

	return nil
}
