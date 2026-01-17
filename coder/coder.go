package coder

import (
	"os"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/parser"
	"github.com/flily/magi-c/preprocessor"
	"github.com/flily/magi-c/tokenizer"
)

const (
	DefaultOutputBase    = "output"
	DefaultMainEntryName = "main"
)

func ParseDocument(data []byte, filename string) (*ast.Document, error) {
	t := tokenizer.NewTokenizerFrom(data, filename)
	parser := parser.NewLLParser(t)
	preprocessor.RegisterPreprocessors(parser)
	return parser.Parse()
}

type Coder struct {
	SourceBase string
	OutputBase string
	Refs       *Cache
}

func NewCoder(sourceBase string, outputBase string) *Coder {
	c := &Coder{
		SourceBase: sourceBase,
		OutputBase: outputBase,
		Refs:       NewCache(),
	}

	return c
}

func (c *Coder) ParseFileContent(filename string, content []byte) error {
	doc, err := ParseDocument(content, filename)
	if err != nil {
		return err
	}

	c.Refs.Add(filename, doc)

	return nil
}

func (c *Coder) ParseFile(filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return c.ParseFileContent(filename, content)
}

func (c *Coder) FindMain() string {
	for filename, doc := range c.Refs.Documents {
		for _, decl := range doc.Declarations {
			if fnDecl, ok := decl.(*ast.FunctionDeclaration); ok {
				if fnDecl.Name.Name == DefaultMainEntryName {
					return filename
				}
			}
		}
	}

	return ""
}
