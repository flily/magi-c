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

	parts := make([]CodeElement, 0, 4+len(v.Declarator)*6)
	parts = append(parts, v.Type)

	for i, decl := range v.Declarator {
		parts = append(parts,
			out.style.Comma().On(i > 0),
			NewElementCollection(
				out.style.PointerSpacingBefore.Select(DelimiterSpace),
				PunctuatorAsterisk.Duplicate(decl.PointerLevel),
				out.style.PointerSpacingAfter.Select(DelimiterSpace),
			).On(decl.PointerLevel > 0),
			NewElementCollection(
				DelimiterSpace,
			).On(decl.PointerLevel <= 0 && i == 0),
			StringElement(decl.Name),
			NewElementCollection(
				out.style.Assign(), decl.Initializer,
			).On(decl.Initializer != nil),
		)
	}

	parts = append(parts, PunctuatorSemicolon, out.style.EOL)
	return out.Write(level, parts...)
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

func (i *ParameterListItem) codeElement() {}

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
	parts := make([]CodeElement, 0, len(p.Items)*2)

	for i, item := range p.Items {
		parts = append(parts, out.style.Comma().On(i > 0), item)
	}

	return out.Write(level, parts...)
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
	err = out.Write(0,
		f.ReturnType, DelimiterSpace, f.Name, OperatorLeftParen, f.Parameters, OperatorRightParen,
		out.style.FunctionNewLine(), OperatorLeftBrace, out.style.EOL,
	)
	if err != nil {
		return err
	}

	for _, stmt := range f.Body {
		if err := stmt.Write(out, level+1); err != nil {
			return err
		}
	}

	if err := out.Write(0, out.style.FunctionBraceIndent, OperatorRightBrace, out.style.EOL); err != nil {
		return err
	}

	return nil
}
