package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/env/incept"
	"github.com/hofstadter-io/hof/lib/runtime"
)

type envFilter func(*env.Env) bool

func commonStart(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole, filters ...envFilter) (R *runtime.Runtime, matches []*env.Env, err error) {
	args, cueargs := splitArgs(args)

	// check the runtime first before starting dagger
	R, err = prepRuntime(cueargs, rflags)
	if err != nil {
		return R, nil, err
	}

	matches = make([]*env.Env, 0)
	for _, e := range R.Envs {
		name, kind, _ := extractMeta(e)
		if name == "" {
			continue
		}

		if !matchValRegexp(name, args) {
			continue
		}
		if !matchValRegexp(kind, eflags.Kind) {
			continue
		}
		if !matchValRegexp(e.Hof.Path, eflags.Path) {
			continue
		}

		if len(filters) > 0 {
			ok := false
			for _, f := range filters {
				if f(e) {
					ok = true
					break
				}
			}

			if !ok {
				continue
			}
		}

		// plan to run the thing
		matches = append(matches, e)
	}

	if len(matches) == 0 {
		return R, nil, fmt.Errorf("no matches found for given args and flags")
	}

	// TODO, sort them somehow

	return R, matches, err
}

func prepRuntime(args []string, rflags flags.RootPflagpole) (*runtime.Runtime, error) {

	// create our core runtime
	r, err := runtime.New(args, rflags)
	if err != nil {
		return nil, err
	}

	err = r.Load()
	if err != nil {
		return nil, cuetils.ExpandCueError(err)
	}

	err = r.EnrichEnv(nil, EnrichEnv)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func EnrichEnv(R *runtime.Runtime, e *env.Env) error {

	// no-op
	return nil
}

func daggerInceptFlags(rflags flags.RootPflagpole, eflags flags.EnvPflagpole) (bool, error) {
	return daggerInceptOpts(&incept.InceptOptions{
		Verbose:     rflags.Verbosity,
		Progress:    eflags.Renderer,
		Interactive: eflags.OnFailure,
		NoExit:      eflags.NoExit,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
		Stdin:       os.Stdin,
	})
}

func daggerInceptOpts(opts *incept.InceptOptions) (bool, error) {
	dst := os.Getenv("DAGGER_SESSION_TOKEN")
	if dst == "" {
		err := incept.Incept(context.Background(), os.Args, opts)
		if err != nil {
			return true, err
		}

		return true, nil
	}
	return false, nil
}

func splitArgs(orig []string) (args, cueargs []string) {
	args = orig
	for i, a := range orig {
		if a == "%" {
			args = orig[:i]
			if i+1 < len(orig) {
				cueargs = orig[i+1:]
			}
			break
		}
	}
	return args, cueargs
}

func extractMeta(e *env.Env) (ename, ekind, mname string) {
	ekind = e.Hof.Env.Kind
	ename = e.Hof.Env.Name
	mname = e.Hof.Metadata.Name
	if ename == "" {
		ename = mname
	}
	return ename, ekind, mname
}

// tries to match val against a list of regexp
// if list is empty, true is returned
func matchValRegexp(val string, res []string) bool {
	if len(res) == 0 {
		return true
	}
	for _, rs := range res {
		re, err := regexp.Compile(rs)
		if err == nil && re.MatchString(val) {
			return true
		}
	}
	return false
}
