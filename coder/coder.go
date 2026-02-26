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
	DefaultOutputBase        = "output"
	DefaultMainEntryName     = "main"
	DefaultOutputSuffix      = ".c"
	DefaultSourceSuffix      = ".mc"
	DefaultOutputParamPrefix = "__out__"
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

func OutputArgumentName(index int) string {
	s := fmt.Sprintf("%s%d", DefaultOutputParamPrefix, index)
	return s
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
	ctx := NewContext()
	elements := c.OutputDeclarations(ctx, document.Declarations)
	return out.Write(0, elements...)
}

func (c *Coder) OutputDeclarations(ctx *Context, decls []ast.Declaration) []csyntax.CodeElement {
	result := make([]csyntax.CodeElement, 0, 2*len(decls))
	for i, decl := range decls {
		out := c.OutputDeclaration(ctx, decl)
		result = append(result, out...)

		if i < len(decls)-1 {
			result = append(result, csyntax.NewEmptyLine())
		}
	}

	return result
}

func (c *Coder) OutputDeclaration(ctx *Context, decl ast.Declaration) []csyntax.CodeElement {
	result := make([]csyntax.CodeElement, 0, 10)
	result = append(result, csyntax.NewContext(decl.Context()))

	switch d := decl.(type) {
	case *ast.FunctionDeclaration:
		result = append(result, c.OutputFunctionDeclaration(ctx, d))

	case *ast.PreprocessorInclude:
		result = append(result, c.OutputPreprocessorInclude(ctx, d))

	case *ast.PreprocessorInline:
		result = append(result, c.OutputPreprocessorInline(ctx, d))
	}

	return result
}

func (c *Coder) OutputFunctionDeclaration(ctx *Context, decl *ast.FunctionDeclaration) *csyntax.FunctionDeclaration {
	if decl.Name.Name == "main" {
		return c.OutputMainFunction(ctx, decl)
	}

	rcc := 0
	if decl.ReturnTypes != nil {
		rcc = decl.ReturnTypes.Length()
	}

	var f *csyntax.FunctionDeclaration
	if rcc > 1 {
		f = c.OutputFunctionMultipleReturnValues(ctx, decl)

	} else {
		f = c.OutputFunctionSingleReturnValue(ctx, decl)
	}

	return f
}

func (c *Coder) outputFunctionBody(ctx *Context, decl *ast.FunctionDeclaration, f *csyntax.FunctionDeclaration) *csyntax.FunctionDeclaration {
	length := len(decl.Statements)
	for i, stmt := range decl.Statements {
		rs := c.OutputStatement(ctx, stmt)
		for _, r := range rs {
			f.AddStatement(r)

		}

		if i < length-1 {
			f.AddStatement(csyntax.NewEmptyLine())
		}
	}

	return f
}

func (c *Coder) OutputMainFunction(ctx *Context, decl *ast.FunctionDeclaration) *csyntax.FunctionDeclaration {
	var f *csyntax.FunctionDeclaration
	if decl.ReturnTypes == nil || decl.ReturnTypes.Length() == 0 {
		f = csyntax.NewFunctionDeclaration("main", csyntax.NewConcreteType("void"), csyntax.NewParameterList(), nil)
	} else {
		f = csyntax.NewFunctionDeclaration("main", csyntax.NewConcreteType("int"), csyntax.NewParameterList(), nil)
	}

	return c.outputFunctionBody(ctx, decl, f)
}

func (c *Coder) OutputFunctionSingleReturnValue(ctx *Context, decl *ast.FunctionDeclaration) *csyntax.FunctionDeclaration {
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
	return c.outputFunctionBody(ctx, decl, f)
}

func (c *Coder) OutputFunctionMultipleReturnValues(ctx *Context, decl *ast.FunctionDeclaration) *csyntax.FunctionDeclaration {
	for i := range decl.ReturnTypes.Types {
		outputParamName := OutputArgumentName(i)
		ctx.FunctionOut.Add(outputParamName, outputParamName)
	}

	retType := csyntax.NewConcreteType("int")
	params := make([]*csyntax.ParameterListItem, 0, 10)
	for _, name := range ctx.FunctionOut.Variables {
		outType := csyntax.NewPointerType("int") // FIXME: output parameter type is always pointer to int for now
		item := csyntax.NewParameterListItem(outType, name.CodeName)
		params = append(params, item)
	}

	if decl.Arguments != nil {
		for _, param := range decl.Arguments.Arguments {
			// FIXME: parameter type is always int for now
			item := csyntax.NewParameterListItem(csyntax.NewConcreteType("int"), param.Name.Name)
			params = append(params, item)
		}
	}

	f := csyntax.NewFunctionDeclaration(decl.Name.Name, retType, csyntax.NewParameterList(params...), nil)
	return c.outputFunctionBody(ctx, decl, f)
}

func (c *Coder) OutputStatement(ctx *Context, stmt ast.Statement) []csyntax.Statement {
	result := make([]csyntax.Statement, 0, 10)
	result = append(result, csyntax.NewContext(stmt.Context()))

	switch s := stmt.(type) {
	case *ast.PreprocessorInclude:
		result = append(result, c.OutputPreprocessorInclude(ctx, s))

	case *ast.PreprocessorInline:
		result = append(result, c.OutputPreprocessorInline(ctx, s))

	case *ast.ReturnStatement:
		result = append(result, c.OutputReturnStatement(ctx, s)...)

	}

	return result
}

func (c *Coder) OutputPreprocessorInclude(ctx *Context, inc *ast.PreprocessorInclude) *csyntax.IncludeDirective {
	var include *csyntax.IncludeDirective

	if inc.LBracket == ast.SLessThan {
		include = csyntax.NewIncludeAngle(inc.Content)

	} else {
		include = csyntax.NewIncludeQuote(inc.Content)
	}

	return include
}

func (c *Coder) OutputPreprocessorInline(ctx *Context, inline *ast.PreprocessorInline) *csyntax.InlineBlock {
	block := csyntax.NewInlineBlock(inline.Content)
	return block
}

func (c *Coder) OutputReturnStatement(ctx *Context, ret *ast.ReturnStatement) []csyntax.Statement {
	stmts := make([]csyntax.Statement, 0, 10)
	if ret.Value == nil || ret.Value.Length() <= 0 {
		stmts = append(stmts, csyntax.NewReturnStatement(nil))
		return stmts
	}

	if ret.Value.Length() == 1 {
		expr := ret.Value.Expressions[0]
		stmt := c.OutputReturnStatementSingleValue(expr.Expression)
		stmts = append(stmts, stmt)
		return stmts
	}

	for i, expr := range ret.Value.Expressions {
		outputParamName := OutputArgumentName(i)
		cexpr := c.OutputExpression(expr.Expression)
		assign := csyntax.NewAssignmentStatement(outputParamName, 1, cexpr) // FIXME: output parameter type is always pointer to concrete type for now
		stmts = append(stmts, assign)
	}

	stmts = append(stmts, csyntax.NewReturnStatement(csyntax.NewIntegerLiteral(0)))
	return stmts
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
