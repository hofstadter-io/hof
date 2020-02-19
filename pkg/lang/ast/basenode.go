package ast

type BaseNode struct {
	ParseInfo *ParseInfo
	parent    ASTNode
}

func (N *BaseNode) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *BaseNode) Parent() ASTNode {
	return N.parent
}
func (N *BaseNode) SetParent(p ASTNode) {
	N.parent = p
}

func (N *BaseNode) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

