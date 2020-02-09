package ast

type TokenPath struct {
	// Parser filled
	ParseInfo *ParseInfo
	Paths []string

	// Phases filled
}

type Token struct {
	// Parser filled
	ParseInfo *ParseInfo
	Value string
}

func (N *TokenPath) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *Token) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

