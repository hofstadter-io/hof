package ast


type Definitions []*Definition

type Definition struct {
	// Parser filled
	ParseInfo *ParseInfo
	Name   *Token
	Target *TokenPath
	Body []ASTNode

	// Phases filled
}

func (N *Definitions) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *Definition) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

