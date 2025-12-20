package csyntax

type VariableDeclaration struct {
	Type         string
	Name         string
	DefaultValue Expression
}

func (v *VariableDeclaration) declarationNode() {}

func (v *VariableDeclaration) statementNode() {}

func (v *VariableDeclaration) Write(out *StyleWriter, level int) error {
	if err := out.WriteIndent(level); err != nil {
		return err
	}

	if err := out.Write("%s %s = ", v.Type, v.Name); err != nil {
		return err
	}

	if err := v.DefaultValue.Write(out, level); err != nil {
		return err
	}

	return out.WriteLine(";")
}
