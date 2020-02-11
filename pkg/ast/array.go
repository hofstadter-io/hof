package ast

type ArrayDef struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode

	Path      *TokenPath

	// Phases filled
}

func (N *ArrayDef) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *ArrayDef) Parent() ASTNode {
	return N.parent
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
	parent    ASTNode

	Elems []ASTNode

	// Phases filled
}

func (N *Array) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *Array) Parent() ASTNode {
	return N.parent
}

func (N *Array) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

