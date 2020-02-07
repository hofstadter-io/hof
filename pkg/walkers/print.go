package walkers

import (
	"errors"
	"fmt"

	"github.com/hofstadter-io/hof/pkg/ast"
)

func Print(node ast.ASTNode) error {
	// fmt.Println("Print BEG")

	switch n := node.(type) {
	case *ast.File:
		// fmt.Println("File:", n.Path)

		err := n.PackageDecl.Walk(Print)
		if err != nil {
			return err
		}

		for _, defn := range n.Definitions {
			err := defn.Walk(Print)
			if err != nil {
				return err
			}
		}

	case *ast.PackageDecl:
		// fmt.Println("Package", n.Name.Value)

	default:
		return errors.New(fmt.Sprintf("%+v", node))
	}

	// fmt.Println("Print END")
	return nil
}
