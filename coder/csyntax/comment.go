package csyntax

type Comment struct {
	Content []string
}

func NewComment(lines ...string) *Comment {
	c := &Comment{
		Content: lines,
	}

	return c
}

func (c *Comment) codeElement()     {}
func (c *Comment) statementNode()   {}
func (c *Comment) declarationNode() {}

func (c *Comment) Write(out *StyleWriter, level int) error {
	if len(c.Content) <= 0 {
		return nil
	}

	indent := out.MakeIndent(level)
	if len(c.Content) == 1 {
		content := StringElement(c.Content[0])
		return out.WriteLine(0, indent, PunctuatorCommentStart, DelimiterSpace, content, DelimiterSpace, PunctuatorCommentEnd)
	}

	parts := make([]CodeElement, 0, 6*len(c.Content)+10)
	// INDENT /* EOL
	parts = append(parts, indent, PunctuatorCommentStart, out.style.EOL)

	for _, line := range c.Content {
		// INDENT SPACE * SPACE line EOL
		parts = append(parts, indent, DelimiterSpace, PunctuatorAsterisk, DelimiterSpace, StringElement(line), out.style.EOL)
	}

	parts = append(parts, indent, DelimiterSpace, PunctuatorCommentEnd, out.style.EOL)
	return out.Write(0, parts...)
}
