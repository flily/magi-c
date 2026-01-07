package csyntax

type IntegerFormat int

const (
	IntegerFormatDecimal IntegerFormat = iota
	IntegerFormatHexadecimalUpper
	IntegerFormatHexadecimalLower
	IntegerFormatOctal
)

type Integer struct {
	Value  int64
	Format IntegerFormat
}

func NewIntegerLiteral(value int64) *Integer {
	i := &Integer{
		Value:  value,
		Format: IntegerFormatDecimal,
	}

	return i
}

func NewHexIntegerLiteralUpper(value int64) *Integer {
	i := &Integer{
		Value:  value,
		Format: IntegerFormatHexadecimalUpper,
	}

	return i
}

func NewHexIntegerLiteralLower(value int64) *Integer {
	i := &Integer{
		Value:  value,
		Format: IntegerFormatHexadecimalLower,
	}

	return i
}

func NewOctalIntegerLiteral(value int64) *Integer {
	i := &Integer{
		Value:  value,
		Format: IntegerFormatOctal,
	}

	return i
}

func (i *Integer) expressionNode() {}

func (i *Integer) codeElement() {}

func (i *Integer) Write(out *StyleWriter, level int) error {
	var elem CodeElement
	switch i.Format {
	case IntegerFormatHexadecimalUpper:
		elem = FormatStringElement("0x%X", i.Value)

	case IntegerFormatHexadecimalLower:
		elem = FormatStringElement("0x%x", i.Value)

	case IntegerFormatOctal:
		elem = FormatStringElement("0%o", i.Value)

	default:
		elem = FormatStringElement("%d", i.Value)
	}

	return out.Write(level, elem)
}
