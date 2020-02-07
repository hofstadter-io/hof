package ast

type Bool struct {
	// Parser filled
	ParseInfo *ParseInfo
	Value bool
}

func (N *Bool) Walk(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}
