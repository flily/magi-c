package csyntax

type FunctionDeclaration struct {
	ReturnType *LiteralType
	Name       StringElement
	Parameters *ParameterList
	Body       *CodeBlock
}

func NewFunctionDeclaration(name string, returnType *LiteralType, parameters *ParameterList, body []Statement) *FunctionDeclaration {
	f := &FunctionDeclaration{
		ReturnType: returnType,
		Name:       StringElement(name),
		Parameters: parameters,
		Body:       NewCodeBlock(body),
	}

	return f
}

func (f *FunctionDeclaration) codeElement()    {}
func (f *FunctionDeclaration) definitionNode() {}

func (f *FunctionDeclaration) AddStatement(stmt Statement) {
	f.Body.Add(stmt)
}

func (f *FunctionDeclaration) Write(out *StyleWriter, level Level) error {
	err := out.WriteIndentLine(level,
		f.ReturnType, DelimiterSpace, f.Name, OperatorLeftParen, f.Parameters, OperatorRightParen,
		out.style.FunctionNewLine(), OperatorLeftBrace, out.style.EOL,
		f.Body,
		out.style.FunctionBraceIndent, OperatorRightBrace)

	return err
}
