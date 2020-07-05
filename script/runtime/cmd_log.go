package runtime

import (
	"fmt"

	"github.com/hofstadter-io/hof/script/ast"
)

func (RT *Runtime) Cmd_log(cmd *ast.Cmd, r *ast.Result) (err error) {
	if cmd.Exp != ast.Pass {
		r.Status = 1
		RT.status = r.Status
		return fmt.Errorf("unsupported: !? exit")
	}

	args := cmd.Args

	if len(args) < 2 {
		r.Status = 1
		RT.status = r.Status
		return fmt.Errorf("usage: log lvl fmt args...")
	}

	lvl, msg, rest := args[0], args[1], args[2:]
	// convert rest for logging
	irest := make([]interface{}, len(rest), len(rest))
	for i := range rest {
		irest[i] = rest[i]
	}

	switch lvl {
	case "debug":
		RT.logger.Debugw(msg, irest...)
	case "debugf":
		RT.logger.Debugf(msg, irest...)

	case "info":
		RT.logger.Infow(msg, irest...)
	case "infof":
		RT.logger.Infof(msg, irest...)

	case "warn":
		RT.logger.Warnw(msg, irest...)
	case "warnf":
		RT.logger.Warnf(msg, irest...)

	case "error":
		RT.logger.Errorw(msg, irest...)
	case "errorf":
		RT.logger.Errorf(msg, irest...)

	case "fatal":
		RT.logger.Fatalw(msg, irest...)
	case "fatalf":
		RT.logger.Fatalf(msg, irest...)

	default:
		return fmt.Errorf("usage: log lvl fmt args...")
	}

	return nil
}
