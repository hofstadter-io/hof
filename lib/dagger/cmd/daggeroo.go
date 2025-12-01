package cmd

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

func Run(id string, rflags flags.RootPflagpole, dflags flags.DaggerooFlagpole) error {
	fmt.Println("veg daggeroo!", id, dflags)
	ctx := context.Background()

	// get dagger client
	dag, err := dagger.Connect(ctx)
	if err != nil {
		return err
	}

	// get directory
	dir := dag.LoadDirectoryFromID(dagger.DirectoryID(id))

	// get container
	run := dag.Container().
		From(dflags.Image).
		WithDirectory("/", dir).
		WithWorkdir(dflags.Workdir).
		WithEnvVariable("CGO_ENABLE", "1")

	_, err = run.Terminal(dagger.ContainerTerminalOpts{Cmd: []string{"/usr/bin/zsh"}}).Sync(ctx)

	return err
}
