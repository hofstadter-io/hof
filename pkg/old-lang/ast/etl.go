package ast

type EtlDefinition struct {
	// Parser filled
	BaseNode

	Name *Token
	Args []*EtlArg
	Return *TokenPath
	Body []ASTNode

	// Phases filled
	SimpleScopedNode
}

type EtlArg struct {
	BaseNode

	Name *Token
	Path *TokenPath
}

