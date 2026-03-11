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

func (b *CodeBlock) codeElement()   {}
func (b *CodeBlock) statementNode() {}

func (b *CodeBlock) Add(stmt Statement) {
	b.Statements = append(b.Statements, stmt)
}

func (b *CodeBlock) Write(out *StyleWriter, level Level) error {
	return out.Write(level.NextIndent(), FromCodeElements(b.Statements...))
}

type EmptyLine struct{}

func NewEmptyLine() *EmptyLine {
	s := &EmptyLine{}

	return s
}

func (s *EmptyLine) codeElement()     {}
func (s *EmptyLine) statementNode()   {}
func (s *EmptyLine) declarationNode() {}

func (s *EmptyLine) Write(out *StyleWriter, level Level) error {
	return out.Write(level, out.style.EOL)
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

func (s *AssignmentStatement) codeElement()   {}
func (s *AssignmentStatement) statementNode() {}

func (s *AssignmentStatement) Write(out *StyleWriter, level Level) error {
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

func (s *ReturnStatement) codeElement()   {}
func (s *ReturnStatement) statementNode() {}

func (s *ReturnStatement) Write(out *StyleWriter, level Level) error {
	return out.WriteIndentLine(level, KeywordReturn,
		NewElementCollection(DelimiterSpace, s.Expression).On(s.Expression != nil),
		PunctuatorSemicolon)
}

type IfStatement struct {
	Expression Expression
	Body       *CodeBlock
}

func NewIfStatement(expression Expression, body *CodeBlock) *IfStatement {
	s := &IfStatement{
		Expression: expression,
		Body:       body,
	}

	return s
}

func (s *IfStatement) codeElement()   {}
func (s *IfStatement) statementNode() {}

func (s *IfStatement) Write(out *StyleWriter, level Level) error {
	parts := []CodeElement{
		KeywordIf, out.style.IfSpacing.Select(DelimiterSpace), OperatorLeftParen, s.Expression, OperatorRightParen,
		out.style.IfNewLine(level), out.style.IfBraceIndent, OperatorLeftBrace, out.style.EOL,
		s.Body,
		out.style.GetIndent(level), out.style.IfBraceIndent, OperatorRightBrace,
	}

	return out.WriteIndentLine(level, parts...)
}

type WhileStatement struct {
	Expression Expression
	Body       *CodeBlock
}

func NewWhileStatement(expression Expression, body *CodeBlock) *WhileStatement {
	s := &WhileStatement{
		Expression: expression,
		Body:       body,
	}

	return s
}

func (s *WhileStatement) codeElement()   {}
func (s *WhileStatement) statementNode() {}

func (s *WhileStatement) Write(out *StyleWriter, level Level) error {
	parts := []CodeElement{
		KeywordWhile, out.style.WhileSpacing.Select(DelimiterSpace), OperatorLeftParen, s.Expression, OperatorRightParen,
		out.style.WhileNewLine(level), out.style.WhileBraceIndent, OperatorLeftBrace, out.style.EOL,
		s.Body,
		out.style.GetIndent(level), out.style.WhileBraceIndent, OperatorRightBrace,
	}

	return out.WriteIndentLine(level, parts...)
}
