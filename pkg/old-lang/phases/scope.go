package phases

import (
	"fmt"
	"reflect"

	"github.com/hofstadter-io/hof/pkg/lang/ast"
	"github.com/hofstadter-io/hof/pkg/lang/context"
)

func FillScopes(ctx *context.Context) error {
	for _, pkg := range ctx.Packages {
		err := fillPackageScope(ctx, pkg)
		if err != nil {
			return err
		}
	}

	return nil
}

func fillPackageScope(ctx *context.Context, pkg *ast.Package) error {
	// fmt.Println("Package:", pkg.Path)
	for _, file := range pkg.Files {
		// fmt.Println(" -", file.Path)
		for _, def := range file.Definitions {

			switch D := def.(type) {
			case *ast.GeneratorDef:
				// fmt.Println("   + gen:", D.Name.Value)
				existing, _ := pkg.LookupInScope([]string{D.Name.Value})
				if existing != nil {
					err := fmt.Errorf("'%s' defined twice in package '%s' %s\n", D.Name.Value, pkg.Path, file.Path)
					ctx.AddError(err)
					continue
				}

				G := &ast.Generator {
					Parsed: D,
					Name: D.Name.Value,
					Paths: D.Paths,
				}

				err := fillGeneratorDefinitionScope(ctx, G)
				if err != nil {
					ctx.AddError(err)
					continue
				}

				err = pkg.DefineInScope(G.Name, G)
				if err != nil {
					ctx.AddError(err)
					continue
				}

			case *ast.TypeDefinition:
				// fmt.Println("   + def:", D.Name.Value)
				existing, _ := pkg.LookupInScope([]string{D.Name.Value})
				if existing != nil {
					err := fmt.Errorf("'%s' defined twice in package '%s' %s\n", D.Name.Value, pkg.Path, file.Path)
					ctx.AddError(err)
					continue
				}

				err := fillTypeDefinitionScope(ctx, D)
				if err != nil {
					ctx.AddError(err)
					continue
				}

				err = pkg.DefineInScope(D.Name.Value, D)
				if err != nil {
					ctx.AddError(err)
					continue
				}

			case *ast.EtlDefinition:
				// fmt.Println("   + def:", D.Name.Value)
				existing, _ := pkg.LookupInScope([]string{D.Name.Value})
				if existing != nil {
					err := fmt.Errorf("'%s' defined twice in package '%s' %s\n", D.Name.Value, pkg.Path, file.Path)
					ctx.AddError(err)
					continue
				}

				err := fillEtlDefinitionScope(ctx, D)
				if err != nil {
					ctx.AddError(err)
					continue
				}

				err = pkg.DefineInScope(D.Name.Value, D)
				if err != nil {
					ctx.AddError(err)
					continue
				}

			default:
				p := reflect.ValueOf(D)
				v := reflect.Indirect(p)
				t := v.Type()
				err := fmt.Errorf("Unknown Definition Type in file '%s:%d'\n%v", file.Path, D.GetParseInfo().Line, t.Name())
				ctx.AddError(err)
				continue

			}
		}
	}

	return nil
}

func fillTypeDefinitionScope(ctx *context.Context, typ *ast.TypeDefinition) error {

	for _, node := range typ.Body {

		switch N := node.(type) {

		case *ast.Field:
			name := N.Key.Value

			p := reflect.ValueOf(N.Value)
			v := reflect.Indirect(p)
			t := v.Type()
			fmt.Printf("%-16s%v\n", name, t.Name())

			typ.DefineInScope(name, N.Value)

		default:
			p := reflect.ValueOf(node)
			v := reflect.Indirect(p)
			t := v.Type()
			fmt.Printf("%-16s%v\n", "unknown field type:", t.Name())

		}
	}

	return nil
}

func fillGeneratorDefinitionScope(ctx *context.Context, gen *ast.Generator) error {

	for _, node := range gen.Parsed.Body {

		switch N := node.(type) {

		case *ast.Field:
			name := N.Key.Value

			p := reflect.ValueOf(N.Value)
			v := reflect.Indirect(p)
			t := v.Type()
			fmt.Printf("%-16s%v\n", name, t.Name())

			gen.DefineInScope(name, N.Value)

		default:
			p := reflect.ValueOf(node)
			v := reflect.Indirect(p)
			t := v.Type()
			fmt.Printf("%-16s%v\n", "unknown field type:", t.Name())

		}
	}

	return nil
}

func fillEtlDefinitionScope(ctx *context.Context, etl *ast.EtlDefinition) error {

	for _, arg := range etl.Args {
		etl.DefineInScope(arg.Name.Value, arg)
	}

	return nil
}
