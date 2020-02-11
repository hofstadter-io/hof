package ast

type Bool struct {
	// Parser filled
	ParseInfo *ParseInfo
	Value bool
}

func (N *Bool) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *Bool) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

type BoolDef struct {
	// Parser filled
	ParseInfo *ParseInfo
}

func (N *BoolDef) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *BoolDef) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}
