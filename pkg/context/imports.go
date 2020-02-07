package context

import (
	// "fmt"

	"github.com/hofstadter-io/hof/pkg/ast"
)

func (ctx *Context) LoadFileImports(file *ast.File) error {
	// fmt.Printf("Load imports for: %#+v\n", file)

	for _, I := range file.Imports {
		path := I.ImportPath.Value

		_, ok := ctx.Packages[path]
		if !ok {
			_, err := ctx.ReadPackage(path, ctx.Module.Config)
			if err != nil {
				ctx.AddError(err)
				ctx.AddPackage(&ast.Package{
					Path: path,
				})
				continue
			}
		}

	}

	return nil
}
