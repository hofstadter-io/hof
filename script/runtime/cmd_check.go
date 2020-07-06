package runtime

import (
	"github.com/hofstadter-io/hof/script/ast"
)

func (RT *Runtime) Cmd_check(cmd *ast.Cmd, r *ast.Result) (err error) {
	RT.logger.Warn("Check", *cmd)

	return nil
}
