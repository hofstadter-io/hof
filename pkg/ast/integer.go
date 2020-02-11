package ast

type IntegerDef struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode
}

func (N *IntegerDef) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *IntegerDef) Parent() ASTNode {
	return N.parent
}

func (N *IntegerDef) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

type Integer struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode

	Value int
}

func (N *Integer) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *Integer) Parent() ASTNode {
	return N.parent
}

func (N *Integer) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

