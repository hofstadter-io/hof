package runtime

import (
	"fmt"
	"strconv"

	"github.com/hofstadter-io/hof/script/ast"
)


// Cmd_status checks the exit or status code from the last effectual command
func (RT *Runtime) Cmd_status(cmd *ast.Cmd, r *ast.Result) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("usage: status <int> ...")
	}

	RT.logger.Debug(cmd, RT.GetStatus())

	// Don't care
	if cmd.Exp == ast.Skip {
		return nil
	}

	// Check arg
	var codes []int
	for i, arg := range cmd.Args {
		code, err := strconv.Atoi(arg)
		if err != nil {
			r.Status = 1
			RT.status = r.Status
			return fmt.Errorf("error: %v\nusage: stdin <int>... (arg %d is not an int)", err, i)
		}
		codes = append(codes, code)
	}

	found := false
	for _, code := range codes {
		if code == RT.GetStatus() {
			found = true
			break
		}
	}

	var err error

	// did we not find what we expected to?
	if cmd.Exp == ast.Pass && !found {
		err = fmt.Errorf("status not found: %d in %v", RT.GetStatus(), codes)
	}

	// did we find something when we expexted not to?
	if cmd.Exp == ast.Fail && found {
		err = fmt.Errorf("unexpected status match: %d in %v", RT.GetStatus(), codes)
	}

	if err != nil {
		r.Status = 1
	} else {
		r.Status = 0
	}
	// XXX Don't update RT.status here, we want to preserve for multiple checks
	// RT.status = r.Status

	return err
}

