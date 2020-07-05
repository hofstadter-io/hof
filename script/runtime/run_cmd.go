package runtime

import (
	"fmt"
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

	// lookup and run command

	// TODO check custom commands

	// check defaults
	C, ok := DefaultCommands[cmd.Cmd]
	if !ok {
		err = fmt.Errorf("Unknown command: %q", cmd.Cmd)
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

