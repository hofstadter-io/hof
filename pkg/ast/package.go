package ast

import (
	"fmt"
	"strings"
	"unicode"
)

type PackageDecl struct {
	// Parser filled
	ParseInfo *ParseInfo
	parent ASTNode
	Name *Token
}

func (N *PackageDecl) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *PackageDecl) Parent() ASTNode {
	return N.parent
}

func (N *PackageDecl) GetParseInfo() *ParseInfo {
	return N.ParseInfo
}

type Package struct {
	Parsed *PackageDecl
	parent ASTNode


	Name string
	// Full import string
	// github.com/hofstadter-io/hof-lang/modules/user
	Path string

	// If defined by a hof-lang.yaml

	// github.com
	Repo string
	// hofstadter-io
	Namespace string
	// hof-lang
	Package string
	// modules/user
	Subpath string

	Files map[string]*File

	Generators map[string]*Generator
	Definitions map[string]*Definition

	publicScope Scope
	privateScope Scope
}

func NewPackage() *Package {
	return &Package {
		Files: map[string]*File{},
		Generators: map[string]*Generator{},
		Definitions: map[string]*Definition{},
		publicScope: map[string]ASTNode{},
		privateScope: map[string]ASTNode{},
	}
}

func (N *Package) GetParseInfo() *ParseInfo {
	return nil
}

func (N *Package) Parent() ASTNode {
	return N.parent
}

func (N *Package) Visit(FN func(ASTNode) (error)) error {
	return FN(N)
}

func (pkg *Package) AddFile(file *File) error {
	path := file.Path
	_, ok := pkg.Files[path]
	if ok {
		// already imported
	} else {
		pkg.Files[path] = file
	}
	return nil
}

func (N *Package) DefineInScope(name string, node ASTNode) error {
	// Check first rune to determine public/private
	// Upper Is Public, lower is private
	r := []rune(name)[0]
	if unicode.IsUpper(r) {
		_, ok := N.publicScope[name]
		if ok {
			return fmt.Errorf("'%s' defined twice", name)
			// return fmt.Errorf("'%s' defined twice\n - %s\n - %s\n", name, existing.GetParseInfo(), node.GetParseInfo())
		}
		N.publicScope[name] = node
	} else {
		_, ok := N.privateScope[name]
		if ok {
			return fmt.Errorf("'%s' defined twice", name)
			// return fmt.Errorf("'%s' defined twice\n - %s\n - %s\n", name, existing.GetParseInfo(), node.GetParseInfo())
		}
		N.privateScope[name] = node
	}
	return nil
}

func (N *Package) LookupInScope(path []string) (ASTNode, error) {
	var err error

	name, rest := path[0], path[1:]
	// Check first rune to determine public/private
	// Upper Is Public, lower is private
	r := []rune(name)[0]
	if unicode.IsUpper(r) {
		existing, ok := N.publicScope[name]
		if ok {
			if len(rest) > 0 {
				return existing.(Scoped).LookupInScope(rest)
			}
			return existing, nil
		}
	} else {
		existing, ok := N.privateScope[name]
		if ok {
			if len(rest) > 0 {
				return existing.(Scoped).LookupInScope(rest)
			}
			return existing, nil
		}
	}

	err = fmt.Errorf("unknown reference to %s", strings.Join(path, "."))
	return nil, err
}

