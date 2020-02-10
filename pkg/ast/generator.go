package ast

type GeneratorDef struct {
	// Parser filled
	ParseInfo *ParseInfo
	Name *Token
	Path *TokenPath
	Extra *Object

	// Phases filled
}

func (N *GeneratorDef) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

type Generator struct {
	Parsed *GeneratorDef
	Name string
	Paths []string
}

func (N *Generator) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

