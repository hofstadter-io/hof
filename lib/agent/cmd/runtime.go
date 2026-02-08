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
	"strconv"
	"strings"

	"cuelang.org/go/cue"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/agent"
	aruntime "github.com/hofstadter-io/hof/lib/agent/runtime"
	"github.com/hofstadter-io/hof/lib/agent/runtime/handlers/ws"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/env/incept"
	"github.com/hofstadter-io/hof/lib/runtime"
)

type agentFilter func(*agent.Agentic) bool

func prepRuntime(args []string, rflags flags.RootPflagpole) (*runtime.Runtime, *aruntime.Runtime, error) {

	// create our core runtime
	r, err := runtime.New(args, rflags)
	if err != nil {
		return nil, nil, err
	}

	err = r.Load()
	if err != nil {
		return nil, nil, cuetils.ExpandCueError(err)
	}

	err = r.InitServices()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init services: %v", err)
	}

	err = r.EnrichAgentic(nil, AgenticEnricher)
	if err != nil {
		return nil, nil, err
	}

	ar, err := aruntime.NewRuntime(r)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create agent runtime: %v", err)
	}

	ar.BackfillAgentic()

	ws.SetupHandlers(r, ar)
	return r, ar, nil
}

func AgenticEnricher(R *runtime.Runtime, e *agent.Agentic) error {
	// no-op for now
	return nil
}

func commonStart(
	args []string,
	rflags flags.RootPflagpole,
	aflags flags.AgentPflagpole,
	filters ...agentFilter,
) (
	R *runtime.Runtime,
	AR *aruntime.Runtime,
	matches []*agent.Agentic,
	err error,
) {
	args, cueargs := cuetils.PercentSplitArgs(args)

	// check the runtime first before starting dagger
	R, AR, err = prepRuntime(cueargs, rflags)
	if err != nil {
		return R, AR, nil, err
	}

	agentics := R.Agentics

	pkg := R.CueConfig.Package
	if pkg == "" {
		pkg = R.BuildInstances[0].PkgName
	}

	if len(rflags.Expression) > 0 {
		disco := make(map[int]*agent.Agentic)
		for _, ex := range rflags.Expression {
			v := cuetils.GetValByEx(ex, pkg, R.Value)
			if v.Exists() {
				switch ik := v.IncompleteKind(); ik {
				case cue.StructKind:
					vp := v.Path().String()
					for i, ev := range agentics {
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
		agentics = slices.Collect(maps.Values(disco))
	}

	// fmt.Println("env.commonStart:", args, pkg, rflags.Expression, len(envs))

	matches = make([]*agent.Agentic, 0)
	for _, e := range agentics {
		// filter for @env() or show all
		if !aflags.ShowAll && !e.Hof.AtMade {
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
		if !matchValRegexp(kind, aflags.Kind) {
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
		return R, AR, nil, fmt.Errorf("no matches found for given args and flags")
	}

	if len(aflags.Sort) > 0 {
		sort.Slice(matches, func(i, j int) bool {
			lhs, rhs := matches[i], matches[j]
			_, lhsKind, lhsMname := extractMeta(lhs)
			_, rhsKind, rhsMname := extractMeta(rhs)

			for _, s := range aflags.Sort {
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

	return R, AR, matches, err
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

func extractMeta(a *agent.Agentic) (aname, akind, mname string) {
	akind = a.Hof.Agentic.Kind
	aname = a.Hof.Agentic.Name
	mname = a.Hof.Metadata.Name
	if aname == "" {
		aname = mname
	}
	if strings.HasPrefix(aname, "\"") && strings.HasSuffix(aname, "\"") {
		aname, _ = strconv.Unquote(aname)
	}
	if strings.HasPrefix(mname, "\"") && strings.HasSuffix(mname, "\"") {
		mname, _ = strconv.Unquote(mname)
	}
	return aname, akind, mname
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
