package csyntax

type Type struct {
	Base         StringElement
	PointerLevel int
}

func NewType(base string, pointerLevel int) *Type {
	t := &Type{
		Base:         StringElement(base),
		PointerLevel: pointerLevel,
	}

	return t
}

func NewConcreteType(base string) *Type {
	return NewType(base, 0)
}

func NewPointerType(base string) *Type {
	return NewType(base, 1)
}

func (t *Type) codeElement() {}

func (t *Type) Write(out *StyleWriter, level int) error {
	parts := []CodeElement{
		t.Base,
		NewElementCollection(
			out.style.PointerSpacingBefore.Select(DelimiterSpace),
			PunctuatorAsterisk.Duplicate(t.PointerLevel),
			out.style.PointerSpacingAfter.Select(DelimiterSpace),
		).On(t.PointerLevel > 0),
	}

	return out.Write(level, parts...)
}

func (t *Type) IsPointer() StyleBoolean {
	return t.PointerLevel > 0
}
