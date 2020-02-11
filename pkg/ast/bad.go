package ast

type BadNode struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode
}

func (N *BadNode) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *BadNode) Parent() ASTNode {
	return N.parent
}

func (N *BadNode) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

