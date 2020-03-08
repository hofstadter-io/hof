package context

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hofstadter-io/hof/pkg/old-lang/ast"
)

func (ctx *Context) LoadFileImports(file *ast.File) error {
	popImportStack := func () {
		ctx.ImportStack = ctx.ImportStack[:len(ctx.ImportStack)-1]
		// fmt.Println("POP-IS", ctx.ImportStack)
	}
	// fmt.Printf("Load imports for: %#+v\n", file)
	if len(ctx.ImportStack) == 0 {
		path := file.Package.Path

		// fix local module packages
		if !strings.HasPrefix(path, "vendor") {
			path = filepath.Join(ctx.Module.Path, path)
			// fmt.Printf("PKG: %#+v\n", *file.Package)
		}

		// Push and defer pop
		ctx.ImportStack = append(ctx.ImportStack, path)
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

		_, ok := ctx.Packages[path]
		if !ok {
			ctx.ImportStack = append(ctx.ImportStack, path)
			_, err := ctx.ReadPackage(path, ctx.Module.Config)
			// don't defer pop here
			popImportStack()

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
