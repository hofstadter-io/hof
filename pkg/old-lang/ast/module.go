package ast

import (
	"github.com/hofstadter-io/hof/pkg/old-lang/module"
)

type Module struct {
	Name string
	// Full import string
	// github.com/hofstadter-io/hof-lang/modules/user
	Path string

	Packages map[string]*Package

	Config *module.Config
}

func (N *Module) GetParseInfo() *ParseInfo {
	return nil
}

func (N *Module) Parent() ASTNode {
	return nil
}

func (N *Module) Visit(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (N *Module) AddPackage(pkg *Package) error {
	path := pkg.Path
	_, ok := N.Packages[path]
	if ok {
		// already imported
	} else {
		N.Packages[path] = pkg
	}
	return nil
}

