package ast

type Definitions []*Definition

func (N *Definitions) GetParseInfo() *ParseInfo {
	return nil
}

func (N *Definitions) Parent() ASTNode {
	return nil
}

func (N *Definitions) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

type Definition struct {
	// Parser filled
	BaseNode

	Name   *Token
	Target *TokenPath
	Body []ASTNode

	// Phases filled
	SimpleScopedNode
}

