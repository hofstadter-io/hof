package runtime

import (
	"github.com/hofstadter-io/hof/script/ast"
)

type RuntimeCommand func (RT *Runtime, cmd *ast.Cmd, result *ast.Result) (err error)

func (RT *Runtime) Cmd_noop(cmd *ast.Cmd, r *ast.Result) (err error) {
	RT.logger.Error("Not Implemented", cmd.Cmd)

	return nil
}


var DefaultCommands = map[string]RuntimeCommand{
	// cmd_call.go
	"call":    (*Runtime).Cmd_noop,

	// cmd_check.go
	"check":    (*Runtime).Cmd_check,

	// cmd_env.go
	"env":     (*Runtime).Cmd_env,

	// cmd_exec.go
	"exec":     (*Runtime).Cmd_exec,

	// cmd_fs.go
	"cd": (*Runtime).Cmd_cd,
	"ls": (*Runtime).Cmd_ls,
	"pwd": (*Runtime).Cmd_pwd,

	"chmod": (*Runtime).Cmd_noop,
	"cp": (*Runtime).Cmd_noop,
	"ln": (*Runtime).Cmd_noop,
	"mkdir": (*Runtime).Cmd_noop,
	"rm": (*Runtime).Cmd_noop,

	// cmd_http.go
	"http": (*Runtime).Cmd_http,

	// cmd_log.go
	"log": (*Runtime).Cmd_log,

	// cmd_status.go
	// TODO, make available to check?
	"status":  (*Runtime).Cmd_status,
	/*
	regexp
	grep
	sed

	also, next section
	*/
	// cmd_stdio.go
	"stdin": (*Runtime).Cmd_noop,
	"stdout": (*Runtime).Cmd_noop,
	"stderr": (*Runtime).Cmd_noop,

	"print": (*Runtime).Cmd_noop,
	"printf": (*Runtime).Cmd_noop,

	/* TODO
	get/set

	sleep
	wait
	skip
	stop
	exit
	sig

	chan?

	alias
	*/
}
