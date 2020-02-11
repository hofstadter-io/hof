package ast

type EtlDefinition struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode

	Name *Token
	Args []*EtlArg
	Return *TokenPath
	Body []ASTNode

	// Phases filled
}

func (N *EtlDefinition) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *EtlDefinition) Parent() ASTNode {
	return N.parent
}

func (N *EtlDefinition) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

type EtlArg struct {
	ParseInfo *ParseInfo
	parent    ASTNode

	Name *Token
	Path *TokenPath
}

func (N *EtlArg) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *EtlArg) Parent() ASTNode {
	return N.parent
}

func (N *EtlArg) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}
