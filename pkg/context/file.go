package context

import (
	"github.com/hofstadter-io/hof/pkg/ast"
	"github.com/hofstadter-io/hof/pkg/config"
	"github.com/hofstadter-io/hof/pkg/parser"
)

func (ctx *Context) ReadFile(filepath string, cfg *config.Config) (*ast.File, error) {

	parseTree, err := parser.ParseHofFile(filepath)
	if err != nil {
		return nil, err
	}

	file := parseTree.(*ast.File)
	file.Path = filepath

	return file, nil
}

