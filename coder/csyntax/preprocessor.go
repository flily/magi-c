package csyntax

type IncludeDirective struct {
	Context
	Filename string
	quoteL   string
	quoteR   string
}

func NewIncludeAngle(ctx *Context, filename string) *IncludeDirective {
	d := &IncludeDirective{
		Context:  *ctx,
		Filename: filename,
		quoteL:   "<",
		quoteR:   ">",
	}

	return d
}

func NewIncludeQuote(ctx *Context, filename string) *IncludeDirective {
	d := &IncludeDirective{
		Context:  *ctx,
		Filename: filename,
		quoteL:   "\"",
		quoteR:   "\"",
	}

	return d
}

func (d *IncludeDirective) Write(out *StyleWriter, level int) error {
	if err := d.Context.Write(out, level); err != nil {
		return err
	}

	return out.WriteLine("#include %s%s%s", d.quoteL, d.Filename, d.quoteR)
}
