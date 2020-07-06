package runtime

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/hofstadter-io/hof/script/ast"
)

func (RT *Runtime) RunNode(node ast.Node, parent *ast.Result) (r *ast.Result, err error) {

	switch node.(type) {

	case *ast.Phase:
		r, err = RT.RunPhase(node.(*ast.Phase), parent)

	case *ast.Cmd:
		r, err = RT.RunCmd(node.(*ast.Cmd), parent)

	case *ast.File:
		r, err = RT.InlineFile(node.(*ast.File), parent)

	default:
		fmt.Printf("  Unhandled Node: %T  %d:%d:%d\n", node, node.DocLine(), node.BegLine(), node.EndLine())
	}

	if err != nil {
		r.AddError(err)
	}

	return r, err
}

func (RT *Runtime) InlineFile(file *ast.File, parent *ast.Result) (r *ast.Result, err error) {
	r = ast.NewResult(file, parent)

	if file.Before {
		return r, nil
	}

	// start result
	r.BegTime = time.Now()
	defer func() {
		if r.EndTime.IsZero() {
			r.EndTime = time.Now()
		}
	}()

	// determine real filename
	filename := RT.MkAbs(file.Path)

	// mkdir if needed
	err = os.MkdirAll(filepath.Dir(filename), 0777)
	if err != nil {
		return r, err
	}

	// write out the file
	err = ioutil.WriteFile(filename, []byte(file.Content), file.Mode)
	if err != nil {
		return r, err
	}

	return r, nil
}
