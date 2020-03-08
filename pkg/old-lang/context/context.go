package context

import (
	// "fmt"

	"github.com/hofstadter-io/hof/pkg/old-lang/ast"
)

type Context struct {
	Entrypoint string
	ImportStack []string

	Module   *ast.Module
	Packages map[string]*ast.Package

	Scope map[string]interface{}
	Errors []error
}

func NewContext() *Context {
	return &Context{
		ImportStack: []string{},
		Packages: map[string]*ast.Package{},
		Scope: map[string]interface{}{},
		Errors: []error{},
	}
}

func (ctx *Context) AddPackage(pkg *ast.Package) error {
	// fmt.Println("AddPackage:", pkg.Path)
	path := pkg.Path
	_, ok := ctx.Packages[path]
	if ok {
		// already imported
	} else {
		ctx.Packages[path] = pkg
	}
	return nil
}

func (ctx *Context) AddError(err error) {
	ctx.Errors = append(ctx.Errors, err)
}

