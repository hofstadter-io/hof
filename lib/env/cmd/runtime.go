package cmd

import (
	"context"
	"fmt"
	"maps"
	"net/url"
	"os"
	"regexp"
	"slices"
	"strings"

	"cuelang.org/go/cue"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/env/incept"
	"github.com/hofstadter-io/hof/lib/runtime"
)

type envFilter func(*env.Env) bool

func isDaggerQueryError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	if !strings.HasPrefix(s, `Post "`) {
		return false
	}
	// find the first quote after Post "
	endIdx := strings.Index(s[6:], `"`)
	if endIdx == -1 {
		return false
	}
	uriStr := s[6 : 6+endIdx]
	u, err := url.Parse(uriStr)
	if err != nil {
		return false
	}
	return u.Path == "/query"
}

func commonStart(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole, filters ...envFilter) (R *runtime.Runtime, matches []*env.Env, err error) {
	args, cueargs := splitArgs(args)

	// check the runtime first before starting dagger
	R, err = prepRuntime(cueargs, rflags)
	if err != nil {
		return R, nil, err
	}

	envs := R.Envs

	pkg := R.CueConfig.Package
	if pkg == "" {
		pkg = R.BuildInstances[0].PkgName
	}

	if len(rflags.Expression) > 0 {
		disco := make(map[int]*env.Env)
		for _, ex := range rflags.Expression {
			v := cuetils.GetValByEx(ex, pkg, R.Value)
			if v.Exists() {
				switch ik := v.IncompleteKind(); ik {
				case cue.StructKind:
					vp := v.Path().String()
					for i, ev := range R.Envs {
						ep := ev.Value.Path().String()
						if strings.HasPrefix(ep, vp) {
							if _, ok := disco[i]; !ok {
								disco[i] = ev
							}
						}
					}

					// This doesn't really work, paths are rewritten by comprehension I think
					// case cue.ListKind:
					// 	iter, err := v.List()
					// 	if err != nil {
					// 		return R, nil, err
					// 	}
					// 	for iter.Next() {
					// 		vp := iter.Value().Path().String()
					// 		for i, ev := range R.Envs {
					// 			ep := ev.Value.Path().String()
					// 			if strings.HasPrefix(ep, vp) {
					// 				if _, ok := disco[i]; !ok {
					// 					disco[i] = ev
					// 				}
					// 			}
					// 		}
					// 	}
				}
			}
		}
		envs = slices.Collect(maps.Values(disco))
	}

	// if len(rflags.Expression) > 0 {
	// 	envs = make([]*env.Env, 0)
	// 	for _, ex := range rflags.Expression {
	// 		// if more than one output, prefix with name in commment
	// 		v := cuetils.GetValByEx(ex, pkg, R.Value)
	// 		if v.Exists() {
	// 			roots, err := hof.FindHofs(v)
	// 			if err != nil {
	// 				return R, nil, err
	// 			}
	// 			for _, root := range roots {
	// 				// what we really want to do is avoid nested nodes
	// 				if root.Value.Path().String() != v.Path().String() {
	// 					continue
	// 				}
	// 				// pkgRoot, pkgPath := root.Value.ReferencePath()
	// 				// fmt.Printf(" - %q %q %q\n", pkgRoot.Path(), pkgPath, root.Value.Path())
	// 				e := &env.Env{Node: root}
	// 				envs = append(envs, e)
	// 			}
	// 		}
	// 	}
	// }

	// this is not right, we need to discover hof.Env within the values returned
	// // if there are -e, build up a list of matching paths
	// if len(rflags.Expression) > 0 {
	// 	envs = make([]*env.Env, 0)
	// 	for _, ex := range rflags.Expression {
	// 		// if more than one output, prefix with name in commment
	// 		v := cuetils.GetValByEx(ex, pkg, R.Value)
	// 		if v.Exists() {
	// 			node, err := hof.ParseHof[any](v)
	// 			if err != nil {
	// 				return R, nil, err
	// 			}
	// 			e := &env.Env{Node: node}
	// 			envs = append(envs, e)
	// 		}
	// 	}
	// }

	// if there are -e, build up a list of matching paths
	// var expaths []string
	// if len(rflags.Expression) > 0 {
	// 	for _, ex := range rflags.Expression {
	// 		// if more than one output, prefix with name in commment
	// 		v := cuetils.GetValByEx(ex, pkg, R.Value)
	// 		if v.Exists() {
	// 			expaths = append(expaths, v.Path().String())
	// 		}
	// 	}
	// }

	// fmt.Println("env.commonStart:", args, pkg, rflags.Expression, len(envs))

	matches = make([]*env.Env, 0)
	for _, e := range envs {
		// filter for @env() or show all
		if !eflags.ShowAll && !e.Hof.AtMade {
			continue
		}
		name, kind, mname := extractMeta(e)
		if name == "" {
			name = mname
		}
		if name == "" {
			continue
		}
		if name == "hide" || name == "hidden" {
			continue
		}

		if !matchValRegexp(mname, args) {
			continue
		}
		if !matchValRegexp(kind, eflags.Kind) {
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

		// do CUE last, since it's the most expensive
		//   we are doing a cheaper, simpler version for now
		// if len(expaths) > 0 {
		// 	p := e.Value.Path().String()
		// 	ok := false
		// 	for _, ex := range expaths {
		// 		if strings.HasPrefix(p, ex) {
		// 			ok = true
		// 			break
		// 		}
		// 	}

		// 	if !ok {
		// 		continue
		// 	}
		// }

		// plan to run the thing
		matches = append(matches, e)
	}

	if len(matches) == 0 {
		return R, nil, fmt.Errorf("no matches found for given args and flags")
	}

	// TODO, sort them somehow for consistency, perhaps moving the list command's logic here

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
