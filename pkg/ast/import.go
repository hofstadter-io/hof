package ast

import (
	"fmt"
	"unicode"
)

type Imports []*Import

func (N *Imports) GetParseInfo() *ParseInfo {
	return nil
}

func (N *Imports) Parent() ASTNode {
	return nil
}

func (N *Imports) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

type Import struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent    ASTNode

	ImportPath   *String
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

func (N *Import) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

func (N *Import) Parent() ASTNode {
	return N.parent
}

func (N *Import) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *Import) DefineInScope(name string, node ASTNode) error {
	return nil
}

func (N *Import) LookupInScope(path []string) (ASTNode, error) {
	name := path[0]
	// Check first rune to determine public/private
	// Upper Is Public, lower is private
	r := []rune(name)[0]
	if unicode.IsUpper(r) {
		return N.Package.LookupInScope(path)
	}

	return nil, fmt.Errorf("Cannot refer to unexported package definitions in '%s'", path)
}
