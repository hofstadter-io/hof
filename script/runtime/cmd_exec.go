package runtime

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hofstadter-io/hof/lib/gotils/intern/os/execpath"
)

// exec runs the given command.
func (ts *Script) CmdExec(neg int, args []string) {

	if len(args) < 1 || (len(args) == 1 && args[0] == "&") {
		ts.Fatalf("usage: exec program [args...] [&]")
	}

	var err error
	if len(args) > 0 && args[len(args)-1] == "&" {
		var cmd *exec.Cmd
		cmd, err = ts.execBackground(args[0], args[1:len(args)-1]...)
		if err == nil {
			wait := make(chan struct{})
			go func() {
				werr := ctxWait(ts.ctxt, cmd)
				close(wait)
				ts.status = cmd.ProcessState.ExitCode()
				err = werr
			}()
			ts.background = append(ts.background, backgroundCmd{cmd, wait, neg})
		}
		ts.stdout, ts.stderr = "", ""
	} else {
		ts.stdout, ts.stderr, err = ts.exec(args[0], args[1:]...)
		if ts.stdout != "" {
			fmt.Fprintf(&ts.log, "[stdout]\n%s", ts.stdout)
		}
		if ts.stderr != "" {
			fmt.Fprintf(&ts.log, "[stderr]\n%s", ts.stderr)
		}
		if err == nil && neg > 0 {
			ts.Fatalf("unexpected command success")
		}
	}

	if err != nil {
		fmt.Fprintf(&ts.log, "[%v]\n", err)
		if ts.ctxt.Err() != nil {
			ts.Fatalf("test timed out while running command")
		} else if neg == 0 {
			ts.Fatalf("unexpected exec command failure")
		}
	}
}


// exec runs the given command line (an actual subprocess, not simulated)
// in ts.cd with environment ts.env and then returns collected standard output and standard error.
func (ts *Script) exec(command string, args ...string) (stdout, stderr string, err error) {
	cmd, err := ts.buildExecCmd(command, args...)
	if err != nil {
		return "", "", err
	}
	cmd.Dir = ts.cd
	cmd.Env = append(ts.env, "PWD="+ts.cd)
	cmd.Stdin = strings.NewReader(ts.stdin)
	var stdoutBuf, stderrBuf strings.Builder
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	if err = cmd.Start(); err == nil {
		err = ctxWait(ts.ctxt, cmd)
		ts.status = cmd.ProcessState.ExitCode()
	}
	ts.stdin = ""
	return stdoutBuf.String(), stderrBuf.String(), err
}

type backgroundCmd struct {
	cmd  *exec.Cmd
	wait <-chan struct{}
	neg  int // if true, cmd should fail
}

// execBackground starts the given command line (an actual subprocess, not simulated)
// in ts.cd with environment ts.env.
func (ts *Script) execBackground(command string, args ...string) (*exec.Cmd, error) {
	cmd, err := ts.buildExecCmd(command, args...)
	if err != nil {
		return nil, err
	}
	cmd.Dir = ts.cd
	cmd.Env = append(ts.env, "PWD="+ts.cd)
	var stdoutBuf, stderrBuf strings.Builder
	cmd.Stdin = strings.NewReader(ts.stdin)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	ts.stdin = ""
	return cmd, cmd.Start()
}

func (ts *Script) buildExecCmd(command string, args ...string) (*exec.Cmd, error) {
	if filepath.Base(command) == command {
		if lp, err := execpath.Look(command, ts.Getenv); err != nil {
			return nil, err
		} else {
			command = lp
		}
	}
	return exec.Command(command, args...), nil
}

// BackgroundCmds returns a slice containing all the commands that have
// been started in the background since the most recent wait command, or
// the start of the script if wait has not been called.
func (ts *Script) BackgroundCmds() []*exec.Cmd {
	cmds := make([]*exec.Cmd, len(ts.background))
	for i, b := range ts.background {
		cmds[i] = b.cmd
	}
	return cmds
}

// ctxWait is like cmd.Wait, but terminates cmd with os.Interrupt if ctx becomes done.
//
// This differs from exec.CommandContext in that it prefers os.Interrupt over os.Kill.
// (See https://golang.org/issue/21135.)
func ctxWait(ctx context.Context, cmd *exec.Cmd) error {
	errc := make(chan error, 1)
	go func() { errc <- cmd.Wait() }()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		interruptProcess(cmd.Process)
		return <-errc
	}
}

// interruptProcess sends os.Interrupt to p if supported, or os.Kill otherwise.
func interruptProcess(p *os.Process) {
	if err := p.Signal(os.Interrupt); err != nil {
		// Per https://golang.org/pkg/os/#Signal, “Interrupt is not implemented on
		// Windows; using it with os.Process.Signal will return an error.”
		// Fall back to Kill instead.
		p.Kill()
	}
}

// Exec runs the given command and saves its stdout and stderr so
// they can be inspected by subsequent script commands.
func (ts *Script) Exec(command string, args ...string) error {
	var err error
	ts.stdout, ts.stderr, err = ts.exec(command, args...)
	if ts.stdout != "" {
		ts.Logf("[stdout]\n%s", ts.stdout)
	}
	if ts.stderr != "" {
		ts.Logf("[stderr]\n%s", ts.stderr)
	}
	return err
}

// Tait waits for background commands to exit, setting stderr and stdout to their result.
func (ts *Script) CmdWait(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? wait")
	}
	if len(args) > 0 {
		ts.Fatalf("usage: wait")
	}

	var stdouts, stderrs []string
	for _, bg := range ts.background {
		<-bg.wait

		args := append([]string{filepath.Base(bg.cmd.Args[0])}, bg.cmd.Args[1:]...)
		fmt.Fprintf(&ts.log, "[background] %s: %v\n", strings.Join(args, " "), bg.cmd.ProcessState)

		cmdStdout := bg.cmd.Stdout.(*strings.Builder).String()
		if cmdStdout != "" {
			fmt.Fprintf(&ts.log, "[stdout]\n%s", cmdStdout)
			stdouts = append(stdouts, cmdStdout)
		}

		cmdStderr := bg.cmd.Stderr.(*strings.Builder).String()
		if cmdStderr != "" {
			fmt.Fprintf(&ts.log, "[stderr]\n%s", cmdStderr)
			stderrs = append(stderrs, cmdStderr)
		}

		if bg.cmd.ProcessState.Success() {
			if bg.neg > 0 {
				ts.Fatalf("unexpected command success")
			}
		} else {
			if ts.ctxt.Err() != nil {
				ts.Fatalf("test timed out while running command")
			} else if bg.neg == 0 {
				ts.Fatalf("unexpected command failure")
			}
		}
	}

	ts.stdout = strings.Join(stdouts, "")
	ts.stderr = strings.Join(stderrs, "")
	ts.background = nil
}

