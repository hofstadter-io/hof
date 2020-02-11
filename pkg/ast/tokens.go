package ast

type Token struct {
	// Parser filled
	ParseInfo *ParseInfo
	Value string
}

func (N *Token) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *Token) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

type TokenPath struct {
	// Parser filled
	ParseInfo *ParseInfo
	Paths []string

	// Phases filled
}

func (N *TokenPath) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *TokenPath) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

