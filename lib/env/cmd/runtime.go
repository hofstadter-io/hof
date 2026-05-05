package cmd

import (
	"context"
	"fmt"
	"maps"
	"net/url"
	"os"
	"regexp"
	"slices"
	"sort"
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
				}
			}
		}
		envs = slices.Collect(maps.Values(disco))
	}

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

		// plan to run the thing
		matches = append(matches, e)
	}

	if len(matches) == 0 {
		return R, nil, fmt.Errorf("no matches found for given args and flags")
	}

	if len(eflags.Sort) > 0 {
		sort.Slice(matches, func(i, j int) bool {
			lhs, rhs := matches[i], matches[j]
			_, lhsKind, lhsMname := extractMeta(lhs)
			_, rhsKind, rhsMname := extractMeta(rhs)

			for _, s := range eflags.Sort {
				s = strings.ToLower(s)
				var lhsVal, rhsVal string
				var desc bool
				s = strings.TrimPrefix(s, "+")
				if c, cut := strings.CutPrefix(s, "-"); cut {
					s = c
					desc = true
				}
				switch s {
				case "name":
					lhsVal, rhsVal = lhsMname, rhsMname
				case "kind":
					lhsVal, rhsVal = lhsKind, rhsKind
				case "path":
					lhsVal, rhsVal = lhs.Hof.Path, rhs.Hof.Path
				case "z":
					if lhs.Hof.Z != rhs.Hof.Z {
						if desc {
							return lhs.Hof.Z > rhs.Hof.Z
						} else {
							return lhs.Hof.Z < rhs.Hof.Z
						}
					}
				default:
					fmt.Printf("WARN: unknown sort field %q ignored", s)
					continue
				}

				if lhsVal != rhsVal {
					if desc {
						return lhsVal > rhsVal
					} else {
						return lhsVal < rhsVal
					}
				}
			}
			return false
		})
	}

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
