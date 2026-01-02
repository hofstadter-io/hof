package cmd

import (
	"fmt"
	"slices"

	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/env/dag"
)

func buildable(e *env.Env) bool {
	accepting := []string{
		"container", "hostImage", "dockerBuild",
		"dir", "hostDir", "gitRepo",
		"file", "hostFile",
	}
	_, kind, _ := extractMeta(e)
	// only publish containers right now
	if slices.Contains(accepting, kind) {
		return true
	}
	return false
}

func Build(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole) error {
	// some quick setup and early filtering
	R, matches, err := commonStart(args, rflags, eflags, buildable)
	if err != nil {
		return err
	}

	// incept if we are not in dagger
	incepted, err := daggerInceptFlags(rflags, eflags)
	if incepted {
		return err
	}

	// setup dagger & cue->dagger engine
	err = R.DaggerInit()
	if err != nil {
		return err
	}
	d, _ := dag.NewClient(R.Ctx, R.DagClient)

	names := make([]string, 0, len(matches))
	for _, e := range matches {
		name, _, _ := extractMeta(e)
		names = append(names, name)
	}

	buildCtx, buildSpan := dagger.Tracer().Start(R.Ctx, fmt.Sprintf("hof env build: %v", names))
	defer buildSpan.End()

	// do actual work
	fmt.Println("building:")
	for i, e := range matches {
		err = func() error {
			name, kind, _ := extractMeta(e)
			matchCtx, matchSpan := dagger.Tracer().Start(buildCtx, fmt.Sprintf("building[%d]: %s (%s)", i, name, kind))
			defer matchSpan.End()
			fmt.Printf(" - %s (%s)\n", name, kind)

			switch kind {
			case "container", "hostImage", "dockerBuild":
				i, err := d.Container(e.Value, eflags.NoCache)
				if err != nil {
					return err
				}
				i, err = i.Sync(matchCtx)
				if err != nil {
					return err
				}

			case "dir", "hostDir", "gitRepo":
				i, _, err := d.Dir(e.Value, eflags.NoCache)
				if err != nil {
					return err
				}
				i, err = i.Sync(matchCtx)
				if err != nil {
					return err
				}

			case "file", "hostFile":
				i, _, err := d.File(e.Value, eflags.NoCache)
				if err != nil {
					return err
				}
				i, err = i.Sync(matchCtx)
				if err != nil {
					return err
				}

			}
			return nil
		}()

		if err != nil {
			return fmt.Errorf("error building match.%d: %v | %v", i, err, e)
		}

	}

	return err
}
