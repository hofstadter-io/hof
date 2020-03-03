package context

import (
	"github.com/hofstadter-io/hof/pkg/lang/ast"
	"github.com/hofstadter-io/hof/pkg/lang/module"
	"github.com/hofstadter-io/hof/pkg/lang/parser"
)

func (ctx *Context) ReadFile(filepath string, cfg *module.Config) (*ast.File, error) {

	parseTree, err := parser.ParseHofFile(filepath)
	if err != nil {
		return nil, err
	}

	file := parseTree.(*ast.File)
	file.Path = filepath

	return file, nil
}

