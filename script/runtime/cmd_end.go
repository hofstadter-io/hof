package runtime

/*
import (
	"github.com/hofstadter-io/hof/script/ast"
)

// skip marks the test skipped.
func (RT *Runtime) Cmd_skip(cmd *ast.Cmd, r *ast.Result) (err error) {
	if neg != 0{
		ts.Fatalf("unsupported: !? skip")
	}

	if len(args) > 1 {
		ts.Fatalf("usage: skip [msg]")
	}

	// Before we mark the test as skipped, shut down any background processes and
	// make sure they have returned the correct status.
	for _, bg := range ts.background {
		interruptProcess(bg.cmd.Process)
	}
	RT.Cmd_wait(0, nil)

	if len(args) == 1 {
		RT.t.Skip(args[0])
	}

	RT.t.Skip()
}

// stop stops execution of the test (marking it passed).
func (RT *Runtime) Cmd_stop(cmd *ast.Cmd, r *ast.Result) (err error) {
	if cmd.Exp != ast.Pass {
		r.Status = 1
		RT.status = r.Status
		return fmt.Errorf("unsupported: !? stop")
	}

	if len(args) > 1 {
		r.Status = 1
		RT.status = r.Status
		return fmt.Errorf("usage: stop [msg]")
	}

	if len(args) == 1 {
		ts.Logf("stop: %s\n", args[0])
	} else {
		ts.Logf("stop\n")
	}

	ts.stopped = true
}

// exit stops execution of the script (marking it according to non-zero exit code).
func (RT *Runtime) Cmd_exit(cmd *ast.Cmd, r *ast.Result) (err error) {
	if cmd.Exp != ast.Pass {
		r.Status = 1
		RT.status = r.Status
		return fmt.Errorf("unsupported: !? exit")
	}

	if len(args) > 2 {
		r.Status = 1
		RT.status = r.Status
		return fmt.Errorf("usage: exit [code] [msg]")
	}

	if len(args) == 2 {
		fmt.Fprintf(cmd.Stdout, "exit: %v %s\n", args[0], args[1])
	} else if len(args) == 1 {
		fmt.Fprintf(cmd.Stdout, "exit: %v %s\n", args[0], args[1])
	} else {
		ts.Logf("exit\n")
	}

	ts.stopped = true
}
*/
