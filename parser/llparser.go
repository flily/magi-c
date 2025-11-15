package parser

import (
	"os"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/preprocessor"
	"github.com/flily/magi-c/tokenizer"
)

type LLParser struct {
	tokenizer  *tokenizer.Tokenizer
	tokens     []ast.Node
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

func (p *LLParser) getToken(index int) ast.Node {
	if index < 0 || index >= len(p.tokens) {
		return nil
	}

	return p.tokens[index]
}

func (p *LLParser) currentToken() ast.Node {
	return p.getToken(p.tokenIndex)
}

func (p *LLParser) peekToken(offset int) ast.Node {
	index := p.tokenIndex + offset
	return p.getToken(index)
}

func (p *LLParser) takeToken() ast.Node {
	token := p.currentToken()
	if token != nil {
		p.tokenIndex++
	}

	return token
}

func (p *LLParser) expectToken(expectedType ast.NodeType) (ast.Node, error) {
	node := p.takeToken()
	if node == nil {
		ctx := p.tokenizer.EOFContext()
		return nil, ast.NewError(ctx, "unexpected EOF, expect token: %s", expectedType.String())
	}

	if node.Type() != expectedType {
		return nil, ast.NewError(node.Context(), "unexpected token: %s, expect: %s", node.Type().String(), expectedType.String())
	}

	return node, nil
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

func (p *LLParser) parseDeclaration(current ast.Node) (ast.Declaration, error) {
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

	rParenArgs, err := p.expectToken(ast.RightParen)
	if err != nil {
		return nil, err
	}
	result.RParenArgs = rParenArgs.(*ast.TerminalToken)

	lParenReturnTypes, err := p.expectToken(ast.LeftParen)
	if err != nil {
		return nil, err
	}
	result.LParenReturnTypes = lParenReturnTypes.(*ast.TerminalToken)

	rParenReturnTypes, err := p.expectToken(ast.RightParen)
	if err != nil {
		return nil, err
	}
	result.RParenReturnTypes = rParenReturnTypes.(*ast.TerminalToken)

	lBracket, err := p.expectToken(ast.LeftBrace)
	if err != nil {
		return nil, err
	}
	result.LBracket = lBracket.(*ast.TerminalToken)

	for {
		current := p.currentToken()
		if current == nil {
			ctx := p.tokenizer.EOFContext()
			return nil, ast.NewError(ctx, "unexpected end of input, expect '}' to close function body")
		}

		if current.Type() == ast.RightBrace {
			result.RBracket = p.takeToken().(*ast.TerminalToken)
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

func (p *LLParser) parseStatement(start ast.Node) (ast.Statement, error) {
	p.takeToken()

	switch start.Type() {
	case ast.Return:
		return p.parseReturn(start.(*ast.TerminalToken))
	}

	return nil, nil
}

func (p *LLParser) parseReturn(keyword *ast.TerminalToken) (ast.Statement, error) {
	result := &ast.ReturnStatement{
		Return: keyword,
	}

	valueNode, err := p.expectToken(ast.Integer)
	if err != nil {
		return nil, err
	}
	result.Value = valueNode.(*ast.IntegerLiteral)

	return result, nil
}
