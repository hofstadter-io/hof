package ast

type Object struct {
	// Parser filled
	ParseInfo *ParseInfo
	Fields []*Field

	// Phases filled
}

func (N *Object) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *Object) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

