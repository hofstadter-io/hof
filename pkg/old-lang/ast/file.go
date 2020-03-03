package ast

import (
	"fmt"
	"os"
	"strings"
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

func (N *File) GetParseInfo() *ParseInfo {
	return nil
}

func (N *File) Parent() ASTNode {
	return N.Package
}

func (N *File) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *File) AddError(err error) {
	N.Errors = append(N.Errors, err)
}

func (N *File) DefineInScope(name string, node ASTNode) error {
	// This will happen during the
	return fmt.Errorf("Files do not have their own scope! You should not use this function")
}

func (N *File) LookupInScope(path []string) (ASTNode, error) {
	var err error
	name, rest := path[0], path[1:]

	// Check Package Scope
	node, err := N.Package.LookupInScope(path)
	if node != nil {
		if len(rest) > 0 {
			return node.(Scoped).LookupInScope(rest)
		} else {
			return node, nil
		}
	}

	// Check Import Scope
	for _, imp := range N.Imports {
		if name == imp.Name {
			return imp.LookupInScope(rest)
		}
	}

	err = fmt.Errorf("unknown reference to %s", strings.Join(path, "."))
	N.AddError(err)
	return nil, err
}

