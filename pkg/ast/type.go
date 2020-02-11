package ast

type TypeDef struct {
	// Parser filled
	ParseInfo *ParseInfo
	Name *Token
	Path *TokenPath
	Extra *Object

	// Phases filled
}

func (N *TypeDef) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *TypeDef) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

