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

func (N *EtlDefinition) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *EtlDefinition) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

func (N *EtlArg) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

type EtlArg struct {
	ParseInfo *ParseInfo
	Name *Token
	Path *TokenPath
}

func (N *EtlArg) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}
