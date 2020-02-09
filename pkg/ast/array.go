package ast

type ArrayDef struct {
	// Parser filled
	ParseInfo *ParseInfo
	Path      *TokenPath

	// Phases filled
}

func (N *ArrayDef) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

type Array struct {
	// Parser filled
	ParseInfo *ParseInfo
	Elems []ASTNode

	// Phases filled
}

func (N *Array) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

