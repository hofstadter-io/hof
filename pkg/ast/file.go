package ast

import (
	"os"
)

type File struct {
	// Parser filled
	PackageDecl *PackageDecl
	Imports []*Import
	Definitions []ASTNode

	// Phases filled
	Errors []error

	Name string
	Path string
	Info os.FileInfo

	Package *Package
}

func (N *File) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}
