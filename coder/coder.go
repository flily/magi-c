package coder

import (
	"fmt"
	"os"
	"path"

	"github.com/flily/magi-c/ast"
	"github.com/flily/magi-c/coder/check"
	"github.com/flily/magi-c/coder/csyntax"
	"github.com/flily/magi-c/context"
	"github.com/flily/magi-c/parser"
	"github.com/flily/magi-c/preprocessor"
	"github.com/flily/magi-c/tokenizer"
)

const (
	DefaultOutputBase    = "output"
	DefaultMainEntryName = "main"
	DefaultOutputSuffix  = ".c"
)

func ParseDocument(data []byte, filename string) (*ast.Document, error) {
	t := tokenizer.NewTokenizerFrom(data, filename)
	parser := parser.NewLLParser(t)
	preprocessor.RegisterPreprocessors(parser)
	return parser.Parse()
}

func OutputFilename(filename string) string {
	return filename + DefaultOutputSuffix
}

type Coder struct {
	SourceBase string
	OutputBase string
	Refs       *Cache
	Style      *csyntax.CodeStyle
}

func NewCoder(sourceBase string, outputBase string) *Coder {
	c := &Coder{
		SourceBase: sourceBase,
		OutputBase: outputBase,
		Refs:       NewCache(),
		Style:      csyntax.KRStyle,
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

func (c *Coder) Check(source string) error {
	doc, ok := c.Refs.Documents[source]
	if !ok {
		return fmt.Errorf("source file '%s' not exists", source)
	}

	conf := check.NewDefaultCheckConfigure()
	checker := check.NewCodeChecker(conf, doc)
	result := checker.Check()
	if result == nil || result.Count(context.Warning) <= 0 {
		return nil
	}

	return result
}

func (c *Coder) Output(source string, target string) error {
	doc, ok := c.Refs.Documents[source]
	if !ok {
		return fmt.Errorf("source file '%s' not exists", source)
	}

	if err := os.MkdirAll(c.OutputBase, 0755); err != nil {
		return err
	}

	fd, err := os.Create(path.Join(c.OutputBase, target))
	if err != nil {
		return err
	}
	defer func() {
		_ = fd.Close()
	}()

	writer := c.Style.MakeWriter(fd)
	return c.OutputDocument(doc, writer)
}

func (c *Coder) OutputDocument(document *ast.Document, out *csyntax.StyleWriter) error {
	decls := make([]csyntax.CodeElement, 0, len(document.Declarations))

	for _, decl := range document.Declarations {
		out := c.OutputDeclaration(decl)
		decls = append(decls, out)
	}

	return out.Write(0, decls...)
}

func (c *Coder) OutputDeclaration(decl ast.Declaration) csyntax.Declaration {
	switch d := decl.(type) {
	case *ast.FunctionDeclaration:
		return c.OutputFunctionDeclaration(d)

	case *ast.PreprocessorInclude:
		return c.OutputPreprocessorInclude(d)

	case *ast.PreprocessorInline:
		return c.OutputPreprocessorInline(d)
	}

	return nil
}

func (c *Coder) OutputFunctionDeclaration(decl *ast.FunctionDeclaration) *csyntax.FunctionDeclaration {
	if decl.Name.Name == "main" {
		return c.OutputMainFunction(decl)
	}

	return nil
}

func (c *Coder) OutputMainFunction(decl *ast.FunctionDeclaration) *csyntax.FunctionDeclaration {
	var f *csyntax.FunctionDeclaration
	if decl.ReturnTypes == nil || decl.ReturnTypes.Length() == 0 {
		f = csyntax.NewFunctionDeclaration("main", csyntax.NewConcreteType("void"), csyntax.NewParameterList(), nil)
	} else {
		f = csyntax.NewFunctionDeclaration("main", csyntax.NewConcreteType("int"), csyntax.NewParameterList(), nil)
	}

	for _, stmt := range decl.Statements {
		r := c.OutputStatement(stmt)
		if r != nil {
			f.AddStatement(r)
		}
	}

	return f
}

func (c *Coder) OutputStatement(stmt ast.Statement) csyntax.Statement {
	switch s := stmt.(type) {
	case *ast.PreprocessorInclude:
		return c.OutputPreprocessorInclude(s)

	case *ast.PreprocessorInline:
		return c.OutputPreprocessorInline(s)

	default:
		return nil
	}
}

func (c *Coder) OutputPreprocessorInclude(inc *ast.PreprocessorInclude) *csyntax.IncludeDirective {
	var include *csyntax.IncludeDirective

	if inc.LBracket == ast.SLessThan {
		include = csyntax.NewIncludeAngle(inc.Context(), inc.Content)

	} else {
		include = csyntax.NewIncludeQuote(inc.Context(), inc.Content)
	}

	return include
}

func (c *Coder) OutputPreprocessorInline(inline *ast.PreprocessorInline) *csyntax.InlineBlock {
	block := csyntax.NewInlineBlock(inline.Context(), inline.Content)
	return block
}
