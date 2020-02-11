package ast

type DecimalDef struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode
}

func (N *DecimalDef) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *DecimalDef) Parent() ASTNode {
	return N.parent
}

func (N *DecimalDef) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

type Decimal struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode

	Value float64
}

func (N *Decimal) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *Decimal) Parent() ASTNode {
	return N.parent
}

func (N *Decimal) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

