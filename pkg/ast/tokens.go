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

func (N *TokenPath) Walk(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *Token) Walk(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

