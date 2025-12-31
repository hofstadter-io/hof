package cmd

import (
	"slices"
	"strings"

	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/env/dag"
)

func runnable(e *env.Env) bool {
	accepting := []string{"container", "hostImage", "dockerBuild"}
	_, kind := extractMeta(e)
	// only publish containers right now
	if slices.Contains(accepting, kind) {
		return true
	}
	return false
}

func Run(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole, cflags flags.Env__RunFlagpole) error {
	// some quick setup and early filtering
	R, matches, err := commonStart(args, rflags, eflags, runnable)
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

	// eventually we want to loop, when we accept more kinds and flags to send them to the background
	e := matches[0]

	i, err := d.Container(e, eflags.NoCache)
	if err != nil {
		return err
	}

	// this terminal thing is ignoring what the container may have set if not built by dagger in this engine
	// do we have a manually set command?
	var cmd []string
	if cflags.Command != "" {
		cmd = strings.Fields(cflags.Command)
	}
	// default args / cmd?
	if len(cmd) == 0 {
		args, _ := i.DefaultArgs(R.Ctx)
		if len(args) > 0 {
			cmd = args
		}
	}
	// entrypoint?
	if len(cmd) == 0 {
		entry, _ := i.Entrypoint(R.Ctx)
		if len(entry) > 0 {
			cmd = entry
		}
	}

	i, err = i.Terminal(dagger.ContainerTerminalOpts{
		Cmd: cmd,
	}).Sync(R.Ctx)
	if err != nil {
		return err
	}

	return nil
}
