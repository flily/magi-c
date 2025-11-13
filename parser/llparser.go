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

func (p *LLParser) currentToken() ast.Node {
	if p.tokenIndex < 0 || p.tokenIndex >= len(p.tokens) {
		return nil
	}

	return p.tokens[p.tokenIndex]
}

func (p *LLParser) takeToken() ast.Node {
	token := p.currentToken()
	if token != nil {
		p.tokenIndex++
	}

	return token
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
	switch current.Type() {
	case ast.NodePreprocessorInclude:
		p.takeToken()
		result = current.(*ast.PreprocessorInclude)

	case ast.NodePreprocessorInline:
		p.takeToken()
		result = current.(*ast.PreprocessorInline)

	default:
		return nil, ast.NewError(current.Context(), "unexpected token: %s", current.Type().String())
	}

	return result, nil
}
