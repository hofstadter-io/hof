package ast

type StringDef struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode
}

func (N *StringDef) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *StringDef) Parent() ASTNode {
	return N.parent
}

func (N *StringDef) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

type String struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode

	Value string
}

func (N *String) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *String) Parent() ASTNode {
	return N.parent
}

func (N *String) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

