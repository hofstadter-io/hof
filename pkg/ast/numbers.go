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

func (N *Integer) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *Decimal) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

type IntegerDef struct {
	// Parser filled
	ParseInfo *ParseInfo
}

func (N *IntegerDef) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

type DecimalDef struct {
	// Parser filled
	ParseInfo *ParseInfo
}

func (N *DecimalDef) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

