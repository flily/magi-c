package csyntax

type VariableDeclarationItem struct {
	PointerLevel int
	Name         string
	Initializer  Expression
}

func NewVariableDeclarator(name string, pointerLevel int, initializer Expression) VariableDeclarationItem {
	d := VariableDeclarationItem{
		PointerLevel: pointerLevel,
		Name:         name,
		Initializer:  initializer,
	}

	return d
}

type VariableDeclaration struct {
	Type       StringElement
	Declarator []VariableDeclarationItem
}

func NewVariableDeclaration(typ string, declarators []VariableDeclarationItem) *VariableDeclaration {
	d := &VariableDeclaration{
		Type:       StringElement(typ),
		Declarator: declarators,
	}

	return d
}

func (v *VariableDeclaration) codeElement()        {}
func (v *VariableDeclaration) declarationNode()    {}
func (v *VariableDeclaration) forInitializerNode() {}

func (v *VariableDeclaration) Write(out *StyleWriter, level Level) error {
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

	return out.Write(level, parts...)
}

func (v *VariableDeclaration) Add(name string, pointerLevel int, initializer Expression) {
	decl := NewVariableDeclarator(name, pointerLevel, initializer)
	v.Declarator = append(v.Declarator, decl)
}

type ParameterListItem struct {
	Type *LiteralType
	Name StringElement
}

func NewParameterListItem(typ *LiteralType, name string) *ParameterListItem {
	p := &ParameterListItem{
		Type: typ,
		Name: StringElement(name),
	}

	return p
}

func (i *ParameterListItem) codeElement() {}

func (i *ParameterListItem) Write(out *StyleWriter, level Level) error {
	return out.Write(level, i.Type, i.Type.IsPointer().Not().Select(DelimiterSpace), i.Name)
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

func (p *ParameterList) Write(out *StyleWriter, level Level) error {
	parts := make([]CodeElement, 0, len(p.Items)*2)

	for i, item := range p.Items {
		parts = append(parts, out.style.Comma().On(i > 0), item)
	}

	return out.Write(level, parts...)
}
