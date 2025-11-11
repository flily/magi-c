package parser

import (
	"os"

	"github.com/flily/magi-c/ast"
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

func (p *LLParser) currentToken() ast.Node {
	if p.tokenIndex < 0 || p.tokenIndex >= len(p.tokens) {
		return nil
	}

	return p.tokens[p.tokenIndex]
}

func (p *LLParser) parseProgram() (*ast.Document, error) {
	statements := make([]ast.Statement, 0, 1000)

	for {
		current := p.currentToken()
		if current == nil {
			break
		}

		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
	}

	document := ast.NewDocument(statements)
	return document, nil
}

func (p *LLParser) parseStatement() (ast.Statement, error) {
	return nil, nil
}
