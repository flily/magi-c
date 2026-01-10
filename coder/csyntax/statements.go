package csyntax

type BlockContext int

const (
	BlockContextFunction BlockContext = iota
	BlockContextIf
	BlockContextElse
	BlockContextDo
	BlockContextWhile
	BlockContextFor
	BlockContextSwitch
	BlockContextCase
	BlockContextDefault
)

type CodeBlock struct {
	Statements []Statement
}

func NewCodeBlock(statements []Statement) *CodeBlock {
	b := &CodeBlock{
		Statements: statements,
	}

	return b
}

func (b *CodeBlock) codeElement() {}

func (b *CodeBlock) statementNode() {}

func (b *CodeBlock) Write(out *StyleWriter, level int) error {
	return out.WriteLine(level, OperatorLeftBrace, out.style.EOL,
		FromCodeElements(b.Statements...),
		out.MakeIndent(level), OperatorRightBrace,
	)
}

type AssignmentStatement struct {
	LeftIdentifier   StringElement
	LeftPointerLevel int
	RightExpression  Expression
}

func NewAssignmentStatement(leftIdentifier string, leftPointerLevel int, rightExpression Expression) *AssignmentStatement {
	s := &AssignmentStatement{
		LeftIdentifier:   StringElement(leftIdentifier),
		LeftPointerLevel: leftPointerLevel,
		RightExpression:  rightExpression,
	}

	return s
}

func (s *AssignmentStatement) codeElement() {}

func (s *AssignmentStatement) statementNode() {}

func (s *AssignmentStatement) Write(out *StyleWriter, level int) error {
	pointer := PunctuatorAsterisk.Duplicate(s.LeftPointerLevel)
	return out.WriteIndentLine(level,
		pointer, out.style.PointerSpacingBefore.Select(DelimiterSpace),
		s.LeftIdentifier, out.style.Assign(), s.RightExpression,
		PunctuatorSemicolon)
}

type ReturnStatement struct {
	Expression Expression
}

func NewReturnStatement(expression Expression) *ReturnStatement {
	s := &ReturnStatement{
		Expression: expression,
	}

	return s
}

func (s *ReturnStatement) statementNode() {}

func (s *ReturnStatement) codeElement() {}

func (s *ReturnStatement) Write(out *StyleWriter, level int) error {
	return out.WriteIndentLine(level, KeywordReturn,
		NewElementCollection(DelimiterSpace, s.Expression).On(s.Expression != nil),
		PunctuatorSemicolon)
}
