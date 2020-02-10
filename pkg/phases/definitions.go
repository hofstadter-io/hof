package phases

import (
	"fmt"

	"github.com/hofstadter-io/hof/pkg/ast"
	"github.com/hofstadter-io/hof/pkg/context"
)

func FindDefinitions(ctx *context.Context) error {
	for _, pkg := range ctx.Packages {
		fmt.Println("Package:", pkg.Path)
    for _, file := range pkg.Files {
      fmt.Println(" -", file.Path)
      for _, def := range file.Definitions {
				switch D := def.(type) {
				case *ast.Generator:
					fmt.Println("   + gen:", D.Name.Value)
					existing, ok := pkg.Generators[D.Name.Value]
					if ok {
						err := fmt.Errorf("Generator '%s' defined twice in package '%s'\n%v\n%v", D.Name.Value, pkg.Path, existing.ParseInfo, D.ParseInfo)
						ctx.AddError(err)
						continue
					}
					pkg.Generators[D.Name.Value] = D

				case *ast.Definition:
					fmt.Println("   + def:", D.Name.Value)
					existing, ok := pkg.Definitions[D.Name.Value]
					if ok {
						err := fmt.Errorf("Definition '%s' defined twice in package '%s'\n%v\n%v", D.Name.Value, pkg.Path, existing.ParseInfo, D.ParseInfo)
						ctx.AddError(err)
						continue
					}
					pkg.Definitions[D.Name.Value] = D
				}
      }
    }
	}

	return nil
}

