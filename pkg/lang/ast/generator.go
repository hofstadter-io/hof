package ast

type GeneratorDef struct {
	// Parser filled
	BaseNode

	Name *Token
	Paths []*TokenPath
	Body []ASTNode
	Extra *Object

	// Phases filled
}

type Generator struct {
	Parsed *GeneratorDef
	BaseNode

	Name string
	Paths []*TokenPath

	PrivacyScopedNode
}

func (N *Generator) GetParseInfo() *ParseInfo {
	return N.Parsed.ParseInfo
}

