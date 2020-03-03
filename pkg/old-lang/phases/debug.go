package phases

import (
  "fmt"
	// "reflect"
	"sort"

	"github.com/hofstadter-io/hof/pkg/lang/ast"
	"github.com/hofstadter-io/hof/pkg/lang/context"
)

func ScopeDebug(ctx *context.Context) error {
	fmt.Println("ScopeDebug\n==========")
	paths := []string{}
	for path, _ := range ctx.Packages {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	for _, path := range paths {
		pkg := ctx.Packages[path]
		keys := pkg.GetScopeKeys()
		fmt.Printf("pkg: %-16s%-64s\n", pkg.Name, pkg.Path)

		for _, key := range keys {
			n, err := pkg.LookupInScope([]string{key})
			if err != nil {
				return err
			}

			/*
			p := reflect.ValueOf(n)
			v := reflect.Indirect(p)
			t := v.Type()
			*/
			// fmt.Printf(" - %-16s%v\n", key, t.Name())

			switch N := n.(type) {
			case *ast.TypeDefinition:
				fmt.Println("  - typ:", key)
				nkeys := N.GetScopeKeys()
				fmt.Printf("    %v\n", nkeys)

			case *ast.Generator:
				fmt.Println("  - gen:", key)
				nkeys := N.GetScopeKeys()
				fmt.Printf("    %v\n", nkeys)

			case *ast.EtlDefinition:
				fmt.Println("  - etl:", key)
				nkeys := N.GetScopeKeys()
				fmt.Printf("    %v\n", nkeys)

			}
		}

		fmt.Println()
	}

	return nil
}


