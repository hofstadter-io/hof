package cmd

import (
	"fmt"
	"slices"

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
	d, _ := dag.NewClient(R.Ctx, R.DagClient)

	// do actual work
	fmt.Println("building:")
	for _, e := range matches {
		name, kind, _ := extractMeta(e)
		fmt.Printf(" - %s (%s)\n", name, kind)

		switch kind {
		case "container", "hostImage", "dockerBuild":
			i, err := d.Container(e, eflags.NoCache)
			if err != nil {
				return err
			}
			i, err = i.Sync(R.Ctx)
			if err != nil {
				return err
			}

		case "dir", "hostDir", "gitRepo":
			i, _, err := d.Dir(e, eflags.NoCache)
			if err != nil {
				return err
			}
			i, err = i.Sync(R.Ctx)
			if err != nil {
				return err
			}

		case "file", "hostFile":
			i, _, err := d.File(e, eflags.NoCache)
			if err != nil {
				return err
			}
			i, err = i.Sync(R.Ctx)
			if err != nil {
				return err
			}

		}

	}

	return err
}
