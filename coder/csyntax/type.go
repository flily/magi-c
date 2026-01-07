package csyntax

import (
	"strings"
)

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
	parts := make([]CodeElement, 0, 10)

	parts = append(parts, t.Base)

	if t.PointerLevel > 0 {
		if out.style.PointerSpacingBefore {
			parts = append(parts, DelimiterSpace)
		}

		pointer := strings.Repeat(PointerAsterisk, t.PointerLevel)
		parts = append(parts, NewStringElement(pointer))

		if out.style.PointerSpacingAfter {
			parts = append(parts, DelimiterSpace)
		}
	}

	return out.WriteItems(level, parts...)
}
