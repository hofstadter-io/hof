package ast

type Integer struct {
	// Parser filled
	ParseInfo *ParseInfo
	Value int
}

type Decimal struct {
	// Parser filled
	ParseInfo *ParseInfo
	Value float64
}

func (N *Integer) Walk(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *Decimal) Walk(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

