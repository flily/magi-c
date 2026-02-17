package coder

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

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
	DefaultSourceSuffix  = ".mc"
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

func (c *Coder) OutputFilename(indexName string) string {
	return path.Join(c.OutputBase, indexName) + DefaultOutputSuffix
}

func (c *Coder) ParseFileContent(filename string, content []byte) (string, error) {
	doc, err := ParseDocument(content, filename)
	if err != nil {
		return "", err
	}

	relName, err := filepath.Rel(c.SourceBase, filename)
	if err != nil {
		panic(err)
	}

	c.Refs.Add(relName, doc)
	return relName, nil
}

func (c *Coder) ParseFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
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

func (c *Coder) Output(sourceRel string) error {
	outputTarget := c.OutputFilename(sourceRel)
	return c.OutputToFile(sourceRel, outputTarget)
}

func (c *Coder) OutputToFile(sourceRel string, target string) error {
	targetBase := path.Dir(target)
	if err := os.MkdirAll(targetBase, 0755); err != nil {
		return err
	}

	fd, err := os.Create(target)
	if err != nil {
		return err
	}
	defer func() {
		_ = fd.Close()
	}()

	return c.OutputTo(sourceRel, fd)
}

func (c *Coder) OutputTo(sourceRel string, out io.StringWriter) error {
	doc, ok := c.Refs.Documents[sourceRel]
	if !ok {
		return fmt.Errorf("source file '%s' not exists", sourceRel)
	}

	writer := c.Style.MakeWriter(out)
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

	rcc := 0
	if decl.ReturnTypes != nil {
		rcc = decl.ReturnTypes.Length()
	}

	var f *csyntax.FunctionDeclaration
	if rcc > 1 {

	} else {
		f = c.OutputFunctionSingleReturnValue(decl)
	}

	return f
}

func (c *Coder) outputFunctionBody(decl *ast.FunctionDeclaration, f *csyntax.FunctionDeclaration) *csyntax.FunctionDeclaration {
	for _, stmt := range decl.Statements {
		r := c.OutputStatement(stmt)
		if r != nil {
			f.AddStatement(r)
		}
	}

	return f
}

func (c *Coder) OutputMainFunction(decl *ast.FunctionDeclaration) *csyntax.FunctionDeclaration {
	var f *csyntax.FunctionDeclaration
	if decl.ReturnTypes == nil || decl.ReturnTypes.Length() == 0 {
		f = csyntax.NewFunctionDeclaration("main", csyntax.NewConcreteType("void"), csyntax.NewParameterList(), nil)
	} else {
		f = csyntax.NewFunctionDeclaration("main", csyntax.NewConcreteType("int"), csyntax.NewParameterList(), nil)
	}

	return c.outputFunctionBody(decl, f)
}

func (c *Coder) OutputFunctionSingleReturnValue(decl *ast.FunctionDeclaration) *csyntax.FunctionDeclaration {
	retType := csyntax.NewConcreteType("void")
	if decl.ReturnTypes != nil && decl.ReturnTypes.Length() > 0 {
		// FIXME: return type is always int for now
		retType = csyntax.NewConcreteType("int")
	}

	params := make([]*csyntax.ParameterListItem, 0, 10)
	if decl.Arguments != nil {
		for _, param := range decl.Arguments.Arguments {
			// FIXME: parameter type is always int for now
			item := csyntax.NewParameterListItem(csyntax.NewConcreteType("int"), param.Name.Name)
			params = append(params, item)
		}
	}

	f := csyntax.NewFunctionDeclaration(decl.Name.Name, retType, csyntax.NewParameterList(params...), nil)
	return c.outputFunctionBody(decl, f)
}

func (c *Coder) OutputStatement(stmt ast.Statement) csyntax.Statement {
	switch s := stmt.(type) {
	case *ast.PreprocessorInclude:
		return c.OutputPreprocessorInclude(s)

	case *ast.PreprocessorInline:
		return c.OutputPreprocessorInline(s)

	case *ast.ReturnStatement:
		return c.OutputReturnStatement(s)

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

func (c *Coder) OutputReturnStatement(ret *ast.ReturnStatement) *csyntax.ReturnStatement {
	if ret.Value == nil || ret.Value.Length() <= 0 {
		return csyntax.NewReturnStatement(nil)
	}

	if ret.Value.Length() == 1 {
		expr := ret.Value.Expressions[0]
		return c.OutputReturnStatementSingleValue(expr.Expression)
	}

	return nil
}

func (c *Coder) OutputReturnStatementSingleValue(expr ast.Expression) *csyntax.ReturnStatement {
	value := c.OutputExpression(expr)
	return csyntax.NewReturnStatement(value)
}

func (c *Coder) OutputExpression(expr ast.Expression) csyntax.Expression {
	switch e := expr.(type) {
	case *ast.Identifier:
		return csyntax.NewIdentifier(e.Name)

	case *ast.IntegerLiteral:
		return csyntax.NewIntegerLiteral(int64(e.Value))

	case *ast.InfixExpression:
		left := c.OutputExpression(e.LeftOperand)
		op := OperatorMap(e.Operator.Token)
		right := c.OutputExpression(e.RightOperand)
		return csyntax.NewInfixExpression(left, op, right)

	default:
		err := fmt.Errorf("unsupported expression type: %T", e)
		panic(err)
	}
}
