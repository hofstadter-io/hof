package context

import (
	"fmt"

	"github.com/hofstadter-io/hof/pkg/ast"
)

func (ctx *Context) LoadFileImports(file *ast.File) error {
	popImportStack := func () {
		ctx.ImportStack = ctx.ImportStack[:len(ctx.ImportStack)-1]
		// fmt.Println("POP-IS", ctx.ImportStack)
	}
	// fmt.Printf("Load imports for: %#+v\n", file)
	if len(ctx.ImportStack) == 0 {
		path := file.Package.Path
		if path == "." {
			path = ctx.Module.Path
			// fmt.Printf("PKG: %#+v\n", *file.Package)
		}
		ctx.ImportStack = append(ctx.ImportStack, path)
		// fmt.Println("PSH-IS", ctx.ImportStack)
		defer popImportStack()
	}

	for _, I := range file.Imports {
		path := I.ImportPath.Value
		for _, imp := range ctx.ImportStack {
			if imp == path {
				werr := fmt.Errorf("Import cycle '%s' in '%s' from %v", path, file.Path, ctx.ImportStack)
				return werr
			}
		}

		ctx.ImportStack = append(ctx.ImportStack, path)
		// fmt.Println("PSH-IS", ctx.ImportStack)
		defer popImportStack()

		_, ok := ctx.Packages[path]
		if !ok {
			_, err := ctx.ReadPackage(path, ctx.Module.Config)
			if err != nil {
				werr := fmt.Errorf("Loading import '%s' in '%s'\n%w", path, file.Path, err)
				ctx.AddError(werr)
				ctx.AddPackage(&ast.Package{
					Path: path,
				})
				continue
			}
		}

	}

	return nil
}
