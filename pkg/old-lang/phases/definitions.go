package phases

import (
	"fmt"
	// "reflect"

	"github.com/hofstadter-io/hof/pkg/old-lang/ast"
	"github.com/hofstadter-io/hof/pkg/old-lang/context"
)

func FindDefinitions(ctx *context.Context) error {
	for _, pkg := range ctx.Packages {
		// fmt.Println("Package:", pkg.Path)
    for _, file := range pkg.Files {
      // fmt.Println(" -", file.Path)
      for _, def := range file.Definitions {

				switch D := def.(type) {
				case *ast.GeneratorDef:
					// fmt.Println("   + gen:", D.Name.Value)
					existing, ok := pkg.Generators[D.Name.Value]
					if ok {
						err := fmt.Errorf("Generator '%s' defined twice in package '%s'\n%v\n%v", D.Name.Value, pkg.Path, existing.Parsed.ParseInfo, D.ParseInfo)
						ctx.AddError(err)
						continue
					}

					G := &ast.Generator {
						Parsed: D,
						Name: D.Name.Value,
						Paths: D.Paths,
					}
					pkg.Generators[G.Name] = G

				case *ast.TypeDefinition:
					// fmt.Println("   + def:", D.Name.Value)
					existing, ok := pkg.Definitions[D.Name.Value]
					if ok {
						err := fmt.Errorf("Definition '%s' defined twice in package '%s'\n%v\n%v", D.Name.Value, pkg.Path, existing.ParseInfo, D.ParseInfo)
						ctx.AddError(err)
						continue
					}
					pkg.Definitions[D.Name.Value] = D

				/*
				default:
					p := reflect.ValueOf(D)
					v := reflect.Indirect(p)
					t := v.Type()
					err := fmt.Errorf("Unknown Definition in file '%s:%d'\n%v", file.Path, D.GetParseInfo().Line, t.Name())
					ctx.AddError(err)
					continue
					*/
				}
      }
    }
	}

	return nil
}

