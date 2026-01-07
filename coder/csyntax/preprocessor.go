package csyntax

type IncludeDirective struct {
	Context
	Filename StringElement
	quoteL   StringElement
	quoteR   StringElement
}

func NewIncludeAngle(ctx *Context, filename string) *IncludeDirective {
	d := &IncludeDirective{
		Context:  *ctx,
		Filename: NewStringElement(filename),
		quoteL:   NewStringElement("<"),
		quoteR:   NewStringElement(">"),
	}

	return d
}

func NewIncludeQuote(ctx *Context, filename string) *IncludeDirective {
	d := &IncludeDirective{
		Context:  *ctx,
		Filename: NewStringElement(filename),
		quoteL:   NewStringElement("\""),
		quoteR:   NewStringElement("\""),
	}

	return d
}

func (d *IncludeDirective) codeElement() {}

func (d *IncludeDirective) declarationNode() {}

func (d *IncludeDirective) statementNode() {}

func (d *IncludeDirective) Write(out *StyleWriter, level int) error {
	if err := d.Context.Write(out, level); err != nil {
		return err
	}

	return out.WriteItems(level, PreprocessorInclude, DelimiterSpace, d.quoteL, d.Filename, d.quoteR, out.EOL)
}
