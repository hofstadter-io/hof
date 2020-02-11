package ast

type GeneratorDef struct {
	// Parser filled
	ParseInfo *ParseInfo
	Name *Token
	Path *TokenPath
	Extra *Object

	// Phases filled
}

func (N *GeneratorDef) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *GeneratorDef) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

type Generator struct {
	Parsed *GeneratorDef
	Name string
	Paths []string
}

func (N *Generator) GetParseInfo() *ParseInfo {
	return N.Parsed.ParseInfo
}

func (N *Generator) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

