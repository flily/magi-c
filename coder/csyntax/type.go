package csyntax

type LiteralType struct {
	Base         StringElement
	PointerLevel int
}

func NewType(base string, pointerLevel int) *LiteralType {
	t := &LiteralType{
		Base:         StringElement(base),
		PointerLevel: pointerLevel,
	}

	return t
}

func NewConcreteType(base string) *LiteralType {
	return NewType(base, 0)
}

func NewPointerType(base string) *LiteralType {
	return NewType(base, 1)
}

func (t *LiteralType) codeElement() {}

func (t *LiteralType) Write(out *StyleWriter, level Level) error {
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

func (t *LiteralType) IsPointer() StyleBoolean {
	return t.PointerLevel > 0
}
