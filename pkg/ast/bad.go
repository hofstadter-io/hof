package ast

type BadNode struct {
	// Parser filled
	ParseInfo *ParseInfo
}

func (N *BadNode) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *BadNode) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

