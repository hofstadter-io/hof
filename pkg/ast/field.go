package ast

type Field struct {
	// Parser filled
	ParseInfo *ParseInfo
	Key   *Token
	Value ASTNode

	// Phases filled
}

func (N *Field) Walk(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

