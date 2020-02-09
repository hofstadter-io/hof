package ast

type String struct {
	// Parser filled
	ParseInfo *ParseInfo
	Value string
}

func (N *String) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

type StringDef struct {
	// Parser filled
	ParseInfo *ParseInfo
}

func (N *StringDef) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

