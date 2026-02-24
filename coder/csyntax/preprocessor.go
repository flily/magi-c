package csyntax

import (
	"github.com/flily/magi-c/context"
)

type IncludeDirective struct {
	Filename StringElement
	quoteL   StringElement
	quoteR   StringElement
}

func NewIncludeAngle(filename string) *IncludeDirective {
	d := &IncludeDirective{
		Filename: StringElement(filename),
		quoteL:   StringElement("<"),
		quoteR:   StringElement(">"),
	}

	return d
}

func NewIncludeQuote(filename string) *IncludeDirective {
	d := &IncludeDirective{
		Filename: StringElement(filename),
		quoteL:   StringElement("\""),
		quoteR:   StringElement("\""),
	}

	return d
}

func (d *IncludeDirective) codeElement()     {}
func (d *IncludeDirective) declarationNode() {}
func (d *IncludeDirective) statementNode()   {}

func (d *IncludeDirective) Write(out *StyleWriter, level int) error {
	return out.WriteLine(level, PreprocessorInclude, DelimiterSpace, d.quoteL, d.Filename, d.quoteR)
}

type InlineBlock struct {
	Context *context.Context
	Content string
}

func NewInlineBlock(content string) *InlineBlock {
	b := &InlineBlock{
		Content: content,
	}

	return b
}

func (b *InlineBlock) codeElement()     {}
func (b *InlineBlock) declarationNode() {}
func (b *InlineBlock) statementNode()   {}

func (b *InlineBlock) Write(out *StyleWriter, level int) error {
	return out.WriteLine(level, StringElement(b.Content))
}
