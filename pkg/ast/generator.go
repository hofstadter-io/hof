package ast

type Generator struct {
	// Parser filled
	ParseInfo *ParseInfo
	Name *Token
	Path *TokenPath
	Extra *Object

	// Phases filled
}

func (N *Generator) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

