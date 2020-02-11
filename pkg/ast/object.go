package ast

type Object struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode

	Fields []*Field

	// Phases filled
	scope Scope
}

func (N *Object) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *Object) Parent() ASTNode {
	return N.parent
}

func (N *Object) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

