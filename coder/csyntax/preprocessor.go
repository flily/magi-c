package csyntax

type IncludeDirective struct {
	Context  *Context
	Filename StringElement
	quoteL   StringElement
	quoteR   StringElement
}

func NewIncludeAngle(ctx *Context, filename string) *IncludeDirective {
	d := &IncludeDirective{
		Context:  ctx,
		Filename: StringElement(filename),
		quoteL:   StringElement("<"),
		quoteR:   StringElement(">"),
	}

	return d
}

func NewIncludeQuote(ctx *Context, filename string) *IncludeDirective {
	d := &IncludeDirective{
		Context:  ctx,
		Filename: StringElement(filename),
		quoteL:   StringElement("\""),
		quoteR:   StringElement("\""),
	}

	return d
}

func (d *IncludeDirective) codeElement() {}

func (d *IncludeDirective) declarationNode() {}

func (d *IncludeDirective) statementNode() {}

func (d *IncludeDirective) Write(out *StyleWriter, level int) error {
	return out.WriteLine(level, d.Context,
		PreprocessorInclude, DelimiterSpace, d.quoteL, d.Filename, d.quoteR)
}
