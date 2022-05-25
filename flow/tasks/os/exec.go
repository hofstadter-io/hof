package os

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"cuelang.org/go/cue"

	hofcontext "github.com/hofstadter-io/hof/flow/context"
)

type Exec struct{}

func NewExec(val cue.Value) (hofcontext.Runner, error) {
	return &Exec{}, nil
}

func (T *Exec) Run(ctx *hofcontext.Context) (interface{}, error) {

	v := ctx.Value
	var cmd *exec.Cmd

	// TODO, rework how i/o works for exec

	// todo, check failure modes, fill, not return error?
	// (in all tasks, really)

	var stdout, stderr io.Writer

	ferr := func() error {
		ctx.CUELock.Lock()
		defer func() {
			ctx.CUELock.Unlock()
		}()
		// get and create command
		cmds, err := extractCmd(v)
		if err != nil {
			return err
		}
		cmd = exec.Command(cmds[0], cmds[1:]...)

		// get dir / env for command
		dir, err := extractDir(v)
		if err != nil {
			return err
		}
		cmd.Dir = dir

		env, err := extractEnv(v)
		if err != nil {
			return err
		}
		cmd.Env = env

		// setup i/o for command
		var stdin io.Reader
		stdin, stdout, stderr, err = extractIO(v)
		if err != nil {
			return err
		}

		if stdin != nil {
			cmd.Stdin = stdin
		}
		if stdout != nil {
			cmd.Stdout = stdout
		}
		if stderr != nil {
			cmd.Stderr = stderr
		}

		return nil
	}()
	if ferr != nil {
		return nil, ferr
	}

	//
	// run command
	//
	err := cmd.Run()

	// TODO, how to run in the background and wait for signal?

	// build return value
	ret := make(map[string]interface{})

	if err != nil {
		ret["error"] = err.Error()
	}

	//
	// possibly fill stdout/stderr
	//
	ret, err = fillIO(v, ret, stdout, stderr)
	if err != nil {
		return nil, err
	}

	// fill exit code / successful
	ret["exitcode"] = cmd.ProcessState.ExitCode()
	ret["success"] = cmd.ProcessState.Success()

	// (TODO): check for user's abort mode preference

	return ret, nil
}

func extractCmd(ex cue.Value) ([]string, error) {
	val := ex.LookupPath(cue.ParsePath("cmd"))
	if val.Err() != nil {
		return nil, val.Err()
	}

	cmds := []string{}
	switch val.IncompleteKind() {
	case cue.StringKind:
		c, err := val.String()
		if err != nil {
			return nil, err
		}
		cmds = []string{c}
	case cue.ListKind:
		l, err := val.List()
		if err != nil {
			return nil, err
		}
		for l.Next() {
			c, err := l.Value().String()
			if err != nil {
				return nil, err
			}
			cmds = append(cmds, c)
		}
	default:
		return nil, fmt.Errorf("unsupported cmd type: ", val.IncompleteKind())
	}

	return cmds, nil
}

func extractDir(ex cue.Value) (string, error) {
	// handle Stdout
	d := ex.LookupPath(cue.ParsePath("dir"))
	if d.Exists() {
		s, err := d.String()
		if err != nil {
			return "", err
		}
		return s, nil
	}
	return "", nil
}

func extractEnv(ex cue.Value) ([]string, error) {

	val := ex.LookupPath(cue.ParsePath("env"))
	if val.Exists() {
		// convert env map in CUE to slice in go
		env := make([]string, 0)
		iter, err := val.Fields()
		if err != nil {
			return nil, err
		}
		for iter.Next() {
			k := iter.Selector().String()
			if err != nil {
				return nil, err
			}
			v, err := iter.Value().String()
			if err != nil {
				return nil, err
			}
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		return env, nil
	}

	return nil, nil
}

func extractIO(ex cue.Value) (Stdin io.Reader, Stdout, Stderr io.Writer, err error) {
	// handle stdin,
	iv := ex.LookupPath(cue.ParsePath("stdin"))
	if iv.Exists() {
		switch iv.IncompleteKind() {
		case cue.StringKind:
			s, err := iv.String()
			if err != nil {
				return nil, nil, nil, err
			}
			if s == "-" {
				// (BUG): works around centralized printing
				Stdin = os.Stdin
			}
			Stdin = strings.NewReader(s)

		case cue.BytesKind:
			b, err := iv.Bytes()
			if err != nil {
				return nil, nil, nil, err
			}
			Stdin = bytes.NewReader(b)

		case cue.NullKind:
			// do nothing so no Stdin is set

		default:
			return nil, nil, nil, fmt.Errorf("unsupported type %v for stdin", iv.IncompleteKind())
		}
	}

	// handle Stdout
	ov := ex.LookupPath(cue.ParsePath("stdout"))
	if !ov.Exists() {
		Stdout = os.Stdout
	} else {
		switch ov.IncompleteKind() {
		// we want a bytes writer for Now
		// will return the proper format when filling the value back
		case cue.StringKind:
			fallthrough
		case cue.BytesKind:
			Stdout = new(bytes.Buffer)

		case cue.NullKind:
			// do nothing so no Stdin is set

		default:
			return nil, nil, nil, fmt.Errorf("unsupported type %v for stdout", ov.IncompleteKind())
		}
	}

	// handle Stderr
	ev := ex.LookupPath(cue.ParsePath("stderr"))
	if !ev.Exists() {
		Stderr = os.Stderr
	} else {
		switch ev.IncompleteKind() {
		// we want a bytes writer for Now
		// will return the proper format when filling the value back
		case cue.StringKind:
			fallthrough
		case cue.BytesKind:
			Stderr = new(bytes.Buffer)

		case cue.NullKind:
			// do nothing so no Stdin is set

		default:
			return nil, nil, nil, fmt.Errorf("unsupported type %v for stderr", ev.IncompleteKind())
		}
	}

	return Stdin, Stdout, Stderr, nil
}

func fillIO(ex cue.Value, ret map[string]interface{}, Stdout, Stderr io.Writer) (map[string]interface{}, error) {
	// (warn) possible cue evaluator race condition here
	ov := ex.LookupPath(cue.ParsePath("stdout"))
	if ov.Exists() {
		switch ov.IncompleteKind() {
		// we want a bytes writer for Now
		// will return the proper format when filling the value back
		case cue.StringKind:
			buf := Stdout.(*bytes.Buffer)
			ret["stdout"] = buf.String()
		case cue.BytesKind:
			buf := Stdout.(*bytes.Buffer)
			ret["stdout"] = buf.Bytes()
		case cue.NullKind:
			// do nothing, Stdout was not captured
		}
	}

	ev := ex.LookupPath(cue.ParsePath("stderr"))
	if ev.Exists() {
		switch ev.IncompleteKind() {
		// we want a bytes writer for Now
		// will return the proper format when filling the value back
		case cue.StringKind:
			buf := Stderr.(*bytes.Buffer)
			ret["stderr"] = buf.String()
		case cue.BytesKind:
			buf := Stderr.(*bytes.Buffer)
			ret["stderr"] = buf.Bytes()
		case cue.NullKind:
			// do nothing, Stderr was not captured
		}
	}

	return ret, nil
}
