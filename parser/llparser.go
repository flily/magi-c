package parser

import (
	"os"
	"slices"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/preprocessor"
	"github.com/flily/magi-c/tokenizer"
)

var expressionFirstSet = []ast.TokenType{
	ast.Integer,
	ast.Float,
	ast.String,
	ast.IdentifierName,
}

func inExpressionFirstSet(t ast.TokenType) bool {
	return slices.Contains(expressionFirstSet, t)
}

type LLParser struct {
	tokenizer  *tokenizer.Tokenizer
	tokens     []ast.TerminalNode
	tokenIndex int
}

func NewLLParser(tokenizer *tokenizer.Tokenizer) *LLParser {
	parser := &LLParser{
		tokenizer:  tokenizer,
		tokens:     nil,
		tokenIndex: 0,
	}

	return parser
}

func NewLLParserFromCode(code string, filename string) *LLParser {
	t := tokenizer.NewTokenizerFromString(code, filename)
	return NewLLParser(t)
}

func NewLLParserFromFile(filename string) (*LLParser, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	t := tokenizer.NewTokenizerFrom(content, filename)
	return NewLLParser(t), nil
}

func (p *LLParser) Parse() (*ast.Document, error) {
	tokens, err := p.tokenizer.ScanAll()
	if err != nil {
		return nil, err
	}

	p.tokens = tokens
	p.tokenIndex = 0
	program, err := p.parseProgram()
	if err != nil {
		return nil, err
	}

	program.Filename = p.tokenizer.Filename
	return program, nil
}

func (p *LLParser) RegisterPreprocessor(name string, handler preprocessor.PreprocessorInitializer) {
	p.tokenizer.RegisterPreprocessor(name, handler)
}

func (p *LLParser) getToken(index int) ast.TerminalNode {
	if index < 0 || index >= len(p.tokens) {
		return nil
	}

	return p.tokens[index]
}

func (p *LLParser) peekToken(offset int) ast.TerminalNode {
	index := p.tokenIndex + offset
	return p.getToken(index)
}

func (p *LLParser) currentToken() ast.TerminalNode {
	return p.peekToken(0)
}

func takeToken[T ast.TerminalNode](p *LLParser) T {
	token := p.currentToken()
	if token != nil {
		p.tokenIndex++
	}

	return token.(T)
}

func (p *LLParser) takeToken() ast.TerminalNode {
	token := p.currentToken()
	if token != nil {
		p.tokenIndex++
	}

	return token
}

func (p *LLParser) restoreToken() ast.TerminalNode {
	if p.tokenIndex > 0 {
		p.tokenIndex--
	}

	return p.currentToken()
}

func (p *LLParser) expectTerminalToken(expectedTypes ...ast.TokenType) (*ast.TerminalToken, error) {
	node := p.currentToken()
	if node == nil {
		ctx := p.tokenizer.EOFContext()
		return nil, ast.NewError(ctx, "unexpected EOF, expect token: %s",
			ast.TokenTypeListString(expectedTypes))
	}

	for _, expectedType := range expectedTypes {
		if node.Type() == expectedType {
			if terminal, ok := node.(*ast.TerminalToken); ok {
				p.takeToken()
				return terminal, nil
			}

			return nil, ast.NewError(node.Context(),
				"unexpected token type %s, expect a terminal token", node.Type())
		}
	}

	return nil, ast.NewError(node.Context(), "unexpected token %s, expect '%s'",
		node.Type(), ast.TokenTypeListString(expectedTypes))
}

func (p *LLParser) expectToken(expectedTypes ...ast.TokenType) (ast.TerminalNode, error) {
	node := p.currentToken()
	if node == nil {
		ctx := p.tokenizer.EOFContext()
		return nil, ast.NewError(ctx, "unexpected EOF, expect token: %s",
			ast.TokenTypeListString(expectedTypes))
	}

	for _, expectedType := range expectedTypes {
		if node.Type() == expectedType {
			p.takeToken()
			return node, nil
		}
	}

	return nil, ast.NewError(node.Context(), "unexpected token %s, expect '%s'",
		node.Type(), ast.TokenTypeListString(expectedTypes))
}

func (p *LLParser) parseProgram() (*ast.Document, error) {
	declarations := make([]ast.Declaration, 0, 1000)

	for {
		current := p.currentToken()
		if current == nil {
			break
		}

		dec, err := p.parseDeclaration(current)
		if err != nil {
			return nil, err
		}

		declarations = append(declarations, dec)
	}

	document := ast.NewDocument(declarations)
	return document, nil
}

func (p *LLParser) parseDeclaration(current ast.TerminalNode) (ast.Declaration, error) {
	var result ast.Declaration
	var err error
	switch current.Type() {
	case ast.NodePreprocessorInclude:
		result = p.takeToken().(*ast.PreprocessorInclude)

	case ast.NodePreprocessorInline:
		result = p.takeToken().(*ast.PreprocessorInline)

	case ast.Function:
		result, err = p.parseFunctionDeclaration()

	default:
		err = ast.NewError(current.Context(), "unexpected token: %s", current.Type().String())
	}

	return result, err
}

func (p *LLParser) parseFunctionDeclaration() (ast.Declaration, error) {
	result := &ast.FunctionDeclaration{
		Keyword: p.takeToken().(*ast.TerminalToken),
	}

	name, err := p.expectToken(ast.IdentifierName)
	if err != nil {
		return nil, err
	}
	result.Name = name.(*ast.Identifier)

	lParenArgs, err := p.expectToken(ast.LeftParen)
	if err != nil {
		return nil, err
	}
	result.LParenArgs = lParenArgs.(*ast.TerminalToken)

	argsLead, err := p.expectToken(ast.RightParen, ast.IdentifierName)
	if err != nil {
		return nil, err
	}

	switch argsLead.Type() {
	case ast.RightParen:
		result.RParenArgs = argsLead.(*ast.TerminalToken)

	case ast.IdentifierName:
		p.restoreToken()
		args, err := p.parseArgumentList()
		if err != nil {
			return nil, err
		}
		result.Arguments = args

		rParenArgs, err := p.expectToken(ast.RightParen)
		if err != nil {
			return nil, err
		}
		result.RParenArgs = rParenArgs.(*ast.TerminalToken)
	}

	lParenReturnTypes, err := p.expectToken(ast.LeftParen)
	if err != nil {
		return nil, err
	}
	result.LParenReturnTypes = lParenReturnTypes.(*ast.TerminalToken)

	typeLead, err := p.expectToken(ast.RightParen, ast.IdentifierName)
	if err != nil {
		return nil, err
	}

	switch typeLead.Type() {
	case ast.RightParen:
		result.RParenReturnTypes = typeLead.(*ast.TerminalToken)

	case ast.IdentifierName:
		p.restoreToken()
		types, err := p.parseTypeList()
		if err != nil {
			return nil, err
		}
		result.ReturnTypes = types

		rParenReturnTypes, err := p.expectToken(ast.RightParen)
		if err != nil {
			return nil, err
		}
		result.RParenReturnTypes = rParenReturnTypes.(*ast.TerminalToken)
	}

	lBracket, err := p.expectToken(ast.LeftBrace)
	if err != nil {
		return nil, err
	}
	result.LBrace = lBracket.(*ast.TerminalToken)

	for {
		current := p.currentToken()
		if current == nil {
			ctx := p.tokenizer.EOFContext()
			return nil, ast.NewError(ctx, "unexpected end of input, expect '}' to close function body")
		}

		if current.Type() == ast.RightBrace {
			result.RBrace = p.takeToken().(*ast.TerminalToken)
			break
		}

		stmt, err := p.parseStatement(current)
		if err != nil {
			return nil, err
		}

		result.Statements = append(result.Statements, stmt)
	}

	return result, nil
}

func (p *LLParser) parseStatement(start ast.TerminalNode) (ast.Statement, error) {
	p.takeToken()

	switch start.Type() {
	case ast.Return:
		return p.parseReturn(start.(*ast.TerminalToken))

	default:
		return nil, ast.NewError(start.Context(), "unexpected token '%s' in statement", start.Type().String())
	}
}

func (p *LLParser) parseReturn(keyword *ast.TerminalToken) (ast.Statement, error) {
	result := ast.NewReturnStatement(keyword)

	returnList, err := p.parseExpressionList()
	if err != nil {
		return nil, err
	}
	result.Value = returnList

	return result, nil
}

func (p *LLParser) parseArgument() (*ast.ArgumentDeclaration, error) {
	arg := ast.NewArgumentDeclaration()
	arg.Name = takeToken[*ast.Identifier](p)

	typeLead := p.currentToken()
	if typeLead == nil {
		ctx := p.tokenizer.EOFContext()
		return nil, ast.NewError(ctx, "unexpected EOF, expect argument type")
	}

	var err error
	var typeNode ast.Type

	switch typeLead.Type() {
	case ast.Asterisk:
		typeNode, err = p.parseSimpleType()

	case ast.IdentifierName:
		typeNode, err = p.parseSimpleType()

	default:
		err = ast.NewError(typeLead.Context(), "unexpected token '%s', expect argument type", typeLead.Type().String())
	}

	if err != nil {
		return nil, err
	}
	arg.Type = typeNode

	last := p.currentToken()
	if last != nil && last.Type() == ast.Comma {
		comma := takeToken[*ast.TerminalToken](p)
		arg.Comma = comma
	}

	return arg, nil
}

func (p *LLParser) parseSimpleType() (*ast.SimpleType, error) {
	t := ast.NewSimpleType()

	for {
		current := p.currentToken()
		if current == nil {
			ctx := p.tokenizer.EOFContext()
			return nil, ast.NewError(ctx, "unexpected EOF, expect type identifier")
		}

		switch current.Type() {
		case ast.Asterisk:
			asterisk := takeToken[*ast.TerminalToken](p)
			t.AddPointerAsterisk(asterisk)

		case ast.IdentifierName:
			identifier := takeToken[*ast.Identifier](p)
			t.Identifier = identifier
			return t, nil

		default:
			return nil, ast.NewError(current.Context(), "unexpected token %s, expect type identifier", current.Type().String())
		}
	}
}

func (p *LLParser) parseArgumentList() (*ast.ArgumentList, error) {
	args := ast.NewArgumentList()

	for {
		current := p.currentToken()
		if current == nil {
			break
		}

		switch current.Type() {
		case ast.RightParen:
			return args, nil

		case ast.IdentifierName:
			arg, err := p.parseArgument()
			if err != nil {
				return nil, err
			}
			args.Arguments = append(args.Arguments, arg)

		default:
			return nil, ast.NewError(current.Context(), "unexpected token '%s' in argument list", current.Type().String())
		}
	}

	return args, nil
}

func (p *LLParser) parseTypeList() (*ast.TypeList, error) {
	types := ast.NewTypeList()

	for {
		current := p.currentToken()
		if current == nil {
			break
		}

		switch current.Type() {
		case ast.RightParen:
			return types, nil

		case ast.IdentifierName, ast.Asterisk:
			typeNode, err := p.parseSimpleType()
			if err != nil {
				return nil, err
			}
			item := ast.NewTypeListItems(typeNode)

			comma := p.currentToken()
			if comma == nil {
				return nil, ast.NewError(p.tokenizer.EOFContext(), "unexpected EOF, expect ',' or ')'")
			}

			if comma.Type() == ast.Comma {
				item.Comma = comma.(*ast.TerminalToken)
				p.takeToken()
			}

			types.Types = append(types.Types, item)

		default:
			return nil, ast.NewError(current.Context(), "unexpected token '%s' in type list", current.Type().String())
		}
	}

	return types, nil
}

func (p *LLParser) parseExpressionList() (*ast.ExpressionList, error) {
	list := ast.NewExpressionList()

	finished := false
	for !finished {
		current := p.currentToken()
		if current == nil {
			break
		}

		if !inExpressionFirstSet(current.Type()) {
			break
		}

		expr, err := p.parseExpression(PrecedenceLowest)
		if err != nil {
			return nil, err
		}

		comma, _ := p.expectTerminalToken(ast.Comma)
		item := ast.NewExpressionListItem(expr, comma)
		list.Expressions = append(list.Expressions, item)
	}

	return nil, nil
}

func (p *LLParser) parseExpression(precedence Precedence) (ast.Expression, error) {
	currrent := p.currentToken()
	if currrent == nil {
		ctx := p.tokenizer.EOFContext()
		return nil, ast.NewError(ctx, "unexpected EOF, expect expression")
	}

	var err error
	var result ast.Expression

	switch currrent.Type() {
	case ast.IdentifierName:
		identifier := takeToken[*ast.Identifier](p)
		result = identifier

	case ast.Integer:
		literal := takeToken[*ast.IntegerLiteral](p)
		result = literal

	default:
		err = ast.NewError(currrent.Context(), "unexpected token '%s' in expression", currrent.Type().String())
	}

	if err != nil {
		return nil, err
	}

	return p.parseComplexExpression(result, precedence)
}

func (p *LLParser) parseComplexExpression(first ast.Expression, precedence Precedence) (ast.Expression, error) {
	current := p.currentToken()
	currentPrecedence := GetPrecedence(current)
	if current == nil || currentPrecedence <= precedence {
		return first, nil
	}

	var expr ast.Expression
	var err error
	switch current.Type() {
	case ast.Plus, ast.Sub:
		expr, err = p.parseInfixExpression(first, precedence)
	}

	if err != nil {
		return nil, err
	}

	return p.parseComplexExpression(expr, precedence)
}

func (p *LLParser) parseInfixExpression(left ast.Expression, precedence Precedence) (*ast.InfixExpression, error) {
	operator := takeToken[*ast.TerminalToken](p)

	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	expr := ast.NewInfixExpression(left, operator, right)
	return expr, nil
}
