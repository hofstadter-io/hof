package cmd

import (
	"fmt"
	"slices"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/env/dag"
)

func publishable(e *env.Env) bool {
	accepting := []string{"container", "hostImage", "dockerBuild"}
	_, kind, _ := extractMeta(e)
	// only publish containers right now
	if slices.Contains(accepting, kind) {
		return true
	}
	return false
}

func Publish(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole, scflags flags.Env__PublishFlagpole) error {
	// some quick setup and early filtering
	R, matches, err := commonStart(args, rflags, eflags, publishable)
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

	fmt.Println("publishing:")
	for _, e := range matches {
		name, kind, _ := extractMeta(e)
		fmt.Printf("  %s (%s)", name, kind)

		i, err := d.Container(e.Value, eflags.NoCache)
		if err != nil {
			return err
		}

		for _, tag := range scflags.Tag {
			uri := fmt.Sprintf("%s/%s:%s", scflags.Registry, name, tag)
			_, err := i.Publish(R.Ctx, uri)
			if err != nil {
				return fmt.Errorf("while publish'n image(%s): %w", uri, err)
			}
			fmt.Println(" -> ", uri)
		}

	}

	return nil
}
