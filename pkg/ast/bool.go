package ast

type BoolDef struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode
}

func (N *BoolDef) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *BoolDef) Parent() ASTNode {
	return N.parent
}

func (N *BoolDef) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

type Bool struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode

	Value bool
}

func (N *Bool) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *Bool) Parent() ASTNode {
	return N.parent
}

func (N *Bool) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}
