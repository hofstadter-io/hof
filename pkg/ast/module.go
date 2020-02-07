package ast

import (
	"github.com/hofstadter-io/hof/pkg/config"
)

type Module struct {
	Name string
	// Full import string
	// github.com/hofstadter-io/hof-lang/modules/user
	Path string

	Packages map[string]*Package

	Config *config.Config
}

func (N *Module) Walk(FN func(ASTNode) (error)) error {
	err := FN(N)
	if err != nil {
		return err
	}

	return nil
}

func (module *Module) AddPackage(pkg *Package) error {
	path := pkg.Path
	_, ok := module.Packages[path]
	if ok {
		// already imported
	} else {
		module.Packages[path] = pkg
	}
	return nil
}

