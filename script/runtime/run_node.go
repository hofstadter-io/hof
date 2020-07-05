package runtime

import (
	"fmt"

	"github.com/hofstadter-io/hof/script/ast"
)

func (RT *Runtime) RunNode(node ast.Node, parent *ast.Result) (r *ast.Result, err error) {

	switch node.(type) {

	case *ast.Phase:
		r, err = RT.RunPhase(node.(*ast.Phase), parent)

	case *ast.Cmd:
		r, err = RT.RunCmd(node.(*ast.Cmd), parent)

	default:
		fmt.Printf("  Unhandled Node: %T  %d:%d:%d\n", node, node.DocLine(), node.BegLine(), node.EndLine())
	}

	if err != nil {
		r.AddError(err)
	}

	return r, err
}
