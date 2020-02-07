package ast

type Imports []*Import

type Import struct {
	// Parser filled
	ParseInfo *ParseInfo
	ImportPath   *Token
	NameOverride *Token

	// Phases filled
	Orig string
	Name string

	Repo string
	Namespace string
	PackageName string
	Subpath string

	Package *Package
}

func (N *Imports) Walk(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *Import) Walk(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

