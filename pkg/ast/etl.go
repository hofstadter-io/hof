package ast

type EtlDefinition struct {
	// Parser filled
	ParseInfo *ParseInfo
	Name *Token
	Args []*EtlArg
	Return *TokenPath
	Body []ASTNode

	// Phases filled
}

func (N *EtlDefinition) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

type EtlArg struct {
	ParseInfo *ParseInfo
	Name *Token
	Path *TokenPath
}

func (N *EtlArg) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}
