package ast

type TypeDef struct {
	// Parser filled
	BaseNode

	Name *Token
	Path *TokenPath
	Extra *Object

	// Phases filled
	PrivacyScopedNode
}

type TypeDefinition struct {
	// Parser filled
	BaseNode

	Open bool

	Name *Token
	Paths []*TokenPath
	Body []ASTNode

	// Phases filled
	SimpleScopedNode
}
