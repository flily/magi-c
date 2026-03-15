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

func (b *CodeBlock) Length() int {
	if b == nil {
		return 0
	}

	return len(b.Statements)
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
	parts := []CodeElement{
		pointer,
		NewElementCollection(
			out.style.PointerSpacingBefore.Select(DelimiterSpace),
		).On(s.LeftPointerLevel > 0),
		s.LeftIdentifier, out.style.Assign(), s.RightExpression,
		PunctuatorSemicolon,
	}

	return out.WriteIndentLine(level, parts...)
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
	ElseBody   *CodeBlock
}

func NewIfElseStatement(expression Expression, thenBody *CodeBlock, elseBody *CodeBlock) *IfStatement {
	s := &IfStatement{
		Expression: expression,
		Body:       thenBody,
		ElseBody:   elseBody,
	}

	return s
}

func NewIfStatement(expression Expression, body *CodeBlock) *IfStatement {
	return NewIfElseStatement(expression, body, nil)
}

func NewIfElseChainStatement(conditions []*IfStatement, elseBody *CodeBlock) *IfStatement {
	var result *IfStatement
	var last *IfStatement

	for _, cond := range conditions {
		if result == nil {
			result = NewIfStatement(cond.Expression, cond.Body)
			last = result

		} else {
			last.ElseBody = NewCodeBlock([]Statement{cond})
			last = cond
		}
	}

	if last != nil {
		last.ElseBody = elseBody
	}

	return result
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

	if s.ElseBody.Length() > 0 {
		first := s.ElseBody.Statements[0]
		if _, ok := first.(*IfStatement); ok {
			parts = append(parts,
				out.style.IfNewLine(level), out.style.IfBraceIndent,
				KeywordElse, DelimiterSpace, first,
			)

			return out.Write(level, parts...)

		} else {
			parts = append(parts,
				out.style.IfNewLine(level), out.style.IfBraceIndent,
				KeywordElse, out.style.IfNewLine(level),
				out.style.IfSpacing.Select(DelimiterSpace), OperatorLeftBrace, out.style.EOL,
				s.ElseBody,
				out.style.GetIndent(level), out.style.IfBraceIndent, OperatorRightBrace,
			)
		}
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

type DoWhileStatement struct {
	Body      *CodeBlock
	Condition Expression
}

func NewDoWhileStatement(body *CodeBlock, condition Expression) *DoWhileStatement {
	s := &DoWhileStatement{
		Body:      body,
		Condition: condition,
	}

	return s
}

func (s *DoWhileStatement) codeElement()   {}
func (s *DoWhileStatement) statementNode() {}

func (s *DoWhileStatement) Write(out *StyleWriter, level Level) error {
	parts := []CodeElement{
		KeywordDo, out.style.DoNewLine(level), out.style.WhileBraceIndent, OperatorLeftBrace, out.style.EOL,
		s.Body,
		out.style.GetIndent(level), out.style.WhileBraceIndent, OperatorRightBrace, DelimiterSpace,
		KeywordWhile, out.style.WhileSpacing.Select(DelimiterSpace), OperatorLeftParen, s.Condition, OperatorRightParen,
		PunctuatorSemicolon,
	}

	return out.WriteIndentLine(level, parts...)
}

type KeywordStatement struct {
	Keyword Keyword
}

func NewKeywordStatement(keyword Keyword) *KeywordStatement {
	s := &KeywordStatement{
		Keyword: keyword,
	}

	return s
}

func NewBreakStatement() *KeywordStatement {
	return NewKeywordStatement(KeywordBreak)
}

func NewContinueStatement() *KeywordStatement {
	return NewKeywordStatement(KeywordContinue)
}

func (s *KeywordStatement) codeElement()   {}
func (s *KeywordStatement) statementNode() {}

func (s *KeywordStatement) Write(out *StyleWriter, level Level) error {
	return out.WriteIndentLine(level, s.Keyword, PunctuatorSemicolon)
}

type CaseBranch struct {
	Expression Expression
	Body       *CodeBlock
}

func NewCaseBranch(expression Expression, body *CodeBlock) *CaseBranch {
	b := &CaseBranch{
		Expression: expression,
		Body:       body,
	}

	return b
}

type SwitchStatement struct {
	Condition Expression
	Cases     []*CaseBranch
	Default   *CodeBlock
}

func NewSwitchStatement(condition Expression, cases []*CaseBranch, defaultBody *CodeBlock) *SwitchStatement {
	s := &SwitchStatement{
		Condition: condition,
		Cases:     cases,
		Default:   defaultBody,
	}

	return s
}

func (s *SwitchStatement) codeElement()   {}
func (s *SwitchStatement) statementNode() {}

func (s *SwitchStatement) Write(out *StyleWriter, level Level) error {
	parts := []CodeElement{
		KeywordSwitch, out.style.SwitchSpacing.Select(DelimiterSpace), OperatorLeftParen, s.Condition, OperatorRightParen,
		out.style.SwitchNewLine(level), out.style.SwitchBraceIndent, OperatorLeftBrace, out.style.EOL,
	}

	for _, caseBranch := range s.Cases {
		parts = append(parts,
			KeywordCase, DelimiterSpace, caseBranch.Expression, PunctuatorColon, out.style.EOL,
			caseBranch.Body,
		)
	}

	if s.Default.Length() > 0 {
		parts = append(parts,
			KeywordDefault, PunctuatorColon, out.style.EOL,
			s.Default,
		)
	}

	parts = append(parts,
		out.style.GetIndent(level), out.style.SwitchBraceIndent, OperatorRightBrace,
	)

	return out.WriteIndentLine(level, parts...)
}
