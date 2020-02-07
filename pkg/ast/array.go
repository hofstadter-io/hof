package ast

type Array struct {
	// Parser filled
	ParseInfo *ParseInfo
	Elems []ASTNode

	// Phases filled
}

func (N *Array) Walk(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

