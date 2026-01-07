package csyntax

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
	Type       StringElement
	Declarator []VariableDeclarator
}

func NewVariableDeclaration(typ string, declarators []VariableDeclarator) *VariableDeclaration {
	d := &VariableDeclaration{
		Type:       StringElement(typ),
		Declarator: declarators,
	}

	return d
}

func (v *VariableDeclaration) codeElement() {}

func (v *VariableDeclaration) declarationNode() {}

func (v *VariableDeclaration) statementNode() {}

func (v *VariableDeclaration) Write(out *StyleWriter, level int) error {
	if err := out.WriteIndent(level); err != nil {
		return err
	}

	if err := out.Write(0, v.Type); err != nil {
		return err
	}

	for i, decl := range v.Declarator {
		if i > 0 {
			if err := out.Write(0, out.style.Comma()); err != nil {
				return err
			}
		}

		pointer := PunctuatorAsterisk.Duplicate(decl.PointerLevel)
		err := out.Write(0,
			NewElementCollection(
				out.style.PointerSpacingBefore.Select(DelimiterSpace),
				pointer,
				out.style.PointerSpacingAfter.Select(DelimiterSpace),
			).On(decl.PointerLevel > 0),
			NewElementCollection(
				DelimiterSpace,
			).On(decl.PointerLevel <= 0 && i == 0),
			StringElement(decl.Name))
		if err != nil {
			return err
		}

		if decl.Initializer != nil {
			if err := out.Write(0, out.style.Assign(), decl.Initializer); err != nil {
				return err
			}
		}
	}

	return out.Write(0, PunctuatorSemicolon, out.EOL)
}

func (v *VariableDeclaration) Add(name string, pointerLevel int, initializer Expression) {
	decl := NewVariableDeclarator(name, pointerLevel, initializer)
	v.Declarator = append(v.Declarator, decl)
}

type ParameterListItem struct {
	Type *Type
	Name StringElement
}

func NewParameterListItem(typ *Type, name string) *ParameterListItem {
	p := &ParameterListItem{
		Type: typ,
		Name: StringElement(name),
	}

	return p
}

func (i *ParameterListItem) Write(out *StyleWriter, level int) error {
	return out.Write(level, i.Type, out.style.PointerSpacingAfter.Select(DelimiterSpace), i.Name)
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

func (p *ParameterList) codeElement() {}

func (p *ParameterList) Write(out *StyleWriter, level int) error {
	for i, item := range p.Items {
		if i > 0 {
			if err := out.Write(0, out.style.Comma()); err != nil {
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
	Name       StringElement
	Parameters *ParameterList
	Body       []Statement
}

func NewFunctionDeclaration(name string, returnType *Type, parameters *ParameterList, body []Statement) *FunctionDeclaration {
	f := &FunctionDeclaration{
		ReturnType: returnType,
		Name:       StringElement(name),
		Parameters: parameters,
		Body:       body,
	}

	return f
}

func (f *FunctionDeclaration) codeElement() {}

func (f *FunctionDeclaration) declarationNode() {}

func (f *FunctionDeclaration) AddStatement(stmt Statement) {
	f.Body = append(f.Body, stmt)
}

func (f *FunctionDeclaration) Write(out *StyleWriter, level int) error {
	if err := out.WriteIndent(level); err != nil {
		return err
	}

	var err error
	err = out.Write(0, f.ReturnType, DelimiterSpace, f.Name, OperatorLeftParen, f.Parameters, OperatorRightParen)
	if err != nil {
		return err
	}

	if out.style.FunctionBraceOnNewLine {
		err = out.Write(0, out.EOL, out.style.FunctionBraceIndent, OperatorLeftBrace, out.EOL)
	} else {
		err = out.Write(0, DelimiterSpace, OperatorLeftBrace, out.EOL)
	}

	for _, stmt := range f.Body {
		if err := stmt.Write(out, level+1); err != nil {
			return err
		}
	}

	if err := out.Write(0, out.style.FunctionBraceIndent, OperatorRightBrace, out.EOL); err != nil {
		return err
	}

	return nil
}
