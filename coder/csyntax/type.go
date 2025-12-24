package csyntax

import (
	"strings"
)

type Type struct {
	Base         string
	PointerLevel int
}

func NewType(base string, pointerLevel int) *Type {
	t := &Type{
		Base:         base,
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

func (t *Type) Write(out *StyleWriter, level int) error {
	if t.PointerLevel <= 0 {
		return out.Write("%s", t.Base)
	}

	format := "%s"
	if out.style.PointerSpacingBefore {
		format = Space + format
	}
	if out.style.PointerSpacingAfter {
		format = format + Space
	}

	pointer := strings.Repeat(PointerAsterisk, t.PointerLevel)
	format = "%s" + format
	return out.Write(format, t.Base, pointer)
}
