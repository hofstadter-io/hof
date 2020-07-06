package runtime

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/hofstadter-io/hof/script/ast"
)

func (RT *Runtime) currCmd(cmd *ast.Cmd) {
	RT.currcmd = cmd
}

func (RT *Runtime) nextCmd(cmd *ast.Cmd) {
	RT.lastcmd = RT.currcmd
	if RT.lastcmd == nil {
		RT.lastcmd = cmd
	}
	RT.currcmd = cmd
}

func (RT *Runtime) RunCmd(cmd *ast.Cmd, parent *ast.Result) (r *ast.Result, err error) {
	// TODO convert to runtime logger
	// fmt.Println("Cmd:", cmd.Cmd, cmd.Args, cmd.DocLine(), cmd.BegLine(), cmd.EndLine())

	// Prep result
	r = ast.NewResult(cmd, parent)
	// TODO, make conditional ?
	// RT.SetMultiWriters(r)

	// start result
	r.BegTime = time.Now()
	defer func() {
		if r.EndTime.IsZero() {
			r.EndTime = time.Now()
		}
	}()

	//
	////// lookup and run command
	//

	found := false

	// check custom commands
	C, ok := RT.params.Cmds[cmd.Cmd]
	if ok {
		found = true
	}

	// check defaults
	if !found {
		C, ok = DefaultCommands[cmd.Cmd]
		if ok {
			found = true
		}
	}

	// check system
	if !found {
		_, err := exec.LookPath(cmd.Cmd)
		if err == nil {
			C, ok = DefaultCommands["exec"]
			if ok {
				found = true
			}
		}
	}

	if !found {
		err = fmt.Errorf("Unknown command: %q", cmd.Cmd)
		r.AddError(err)
		return r, err
	}

	fmt.Fprintln(RT.Stdout, ">>>", cmd)
	err = C(RT, cmd, r)
	if err != nil {
		r.AddError(err)
	}
	r.EndTime = time.Now()

	fmt.Fprintf(RT.Stdout, "[status:%d] [time:%v]\n", r.Status, r.EndTime.Sub(r.BegTime))

	stdout := r.Stdout.(fmt.Stringer).String()
	if stdout != "" {
		fmt.Fprintf(RT.Stdout, "[stdout]\n%s\n", stdout)
	}

	stderr := r.Stderr.(fmt.Stringer).String()
	if stderr != "" {
		fmt.Fprintf(RT.Stderr, "[stderr]\n%s\n", stderr)
	}

	/*
	if cmd.Exp == ast.Pass && len(r.Errors) > 0 {
		return r, fmt.Errorf("expected cmd to pass: %d", r.Status)
	}

	if cmd.Exp == ast.Fail && len(r.Errors) == 0 {
		return r, fmt.Errorf("expected cmd to fail: %d", r.Status)
	}
	*/

	if cmd.Exp != ast.Skip && len(r.Errors) > 0 {
		for _, e := range r.Errors {
			RT.logger.Error(e)
		}
		return r, fmt.Errorf("%d errors occurred", len(r.Errors))
	}

	return r, nil
}

