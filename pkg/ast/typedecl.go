package ast

type TypeDecl struct {
	// Parser filled
	ParseInfo *ParseInfo
	Name *Token
	Path *TokenPath
	Extra *Object

	// Phases filled
}

func (N *TypeDecl) Walk(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

