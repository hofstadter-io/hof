package cmd

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"dagger.io/dagger"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/env/dag"
	"golang.org/x/sync/errgroup"
)

var accepting = []string{
	"container",
	"dockerBuild",
	"service",

	"dir",
	"file",
	"secret",

	"gitRepo",
	"hostDir",
	"hostFile",
	"hostImage",
	// "hostService",
	// "hostTunnel",
	// "hostSocket",

	"exportDir",
	"exportFile",
	"exportImageFile",
	"exportImage",
	"publishImage",

	"cuefigSBOM",
}

func syncable(e *env.Env) bool {
	_, kind, _ := extractMeta(e)
	// only publish containers right now
	if slices.Contains(accepting, kind) {
		return true
	}
	return false
}

type Syncable[T any] interface {
	Sync(context.Context) (T, error)
}

func maybeSync[T Syncable[T]](ctx context.Context, val T, dryRun bool) error {
	if dryRun {
		return nil
	}
	_, err := val.Sync(ctx)
	return err
}

func Sync(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole) error {
	veryStart := time.Now()
	// some quick setup and early filtering
	R, matches, err := commonStart(args, rflags, eflags, syncable)
	if err != nil {
		return err
	}

	// incept if we are not in dagger
	incepted, err := daggerInceptFlags(rflags, eflags)
	if incepted {
		return err
	}
	defer func() {
		fmt.Println("done! ", time.Since(veryStart).Round(time.Millisecond))
	}()
	fmt.Printf("init'n  ")

	// setup dagger & cue->dagger engine
	err = R.DaggerInit()
	if err != nil {
		return err
	}
	d, _ := dag.NewClient(R.Ctx, R.DAG)

	buildCtx, buildSpan := dagger.Tracer().Start(R.Ctx, "hof env sync")
	defer buildSpan.End()

	// A helper function to handle the common DryRun + Sync logic

	// do actual work
	fmt.Printf("  %v\n", time.Since(veryStart).Round(time.Millisecond))
	fmt.Println("sync'n")

	var g *errgroup.Group
	var groupCtx context.Context
	if eflags.FailFast {
		g, groupCtx = errgroup.WithContext(buildCtx)
	} else {
		g = new(errgroup.Group)
		groupCtx = buildCtx
	}

	if eflags.Parallel > 0 {
		g.SetLimit(eflags.Parallel)
	}

	for i, e := range matches {
		i, e := i, e
		g.Go(func() error {
			name, kind, _ := extractMeta(e)
			matchCtx, matchSpan := dagger.Tracer().Start(groupCtx, fmt.Sprintf("building[%d]: %s (%s)", i, name, kind))
			defer matchSpan.End()

			if groupCtx.Err() != nil {
				fmt.Printf("ABORT: %s (%s)\n", name, kind)
				return groupCtx.Err()
			}

			fmt.Printf("START: %s (%s)\n", name, kind)
			start := time.Now()

			var err error
			// it would be freaking sweet if there was a way to get generics or something to move these syncs out and have just one, but different types and returns

			switch kind {
			case "container", "dockerBuild", "hostImage":
				val, err2 := d.Container(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				err = maybeSync(matchCtx, val, rflags.DryRun)

			case "dir", "gitRepo", "hostDir":
				val, _, err2 := d.Dir(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				err = maybeSync(matchCtx, val, rflags.DryRun)

			case "file", "hostFile":
				val, _, err2 := d.File(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				err = maybeSync(matchCtx, val, rflags.DryRun)

			case "exportDir":
				val, _, err2 := d.HashExportDir(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				if rflags.DryRun {
					break
				}
				_, err = val.Sync(matchCtx)

			case "exportFile":
				val, _, err2 := d.HashExportFile(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				err = maybeSync(matchCtx, val, rflags.DryRun)

			case "exportImageFile":
				val, _, err2 := d.HashExportImageFile(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				err = maybeSync(matchCtx, val, rflags.DryRun)

			case "exportImage":
				val, _, err2 := d.HashExportImage(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				err = maybeSync(matchCtx, val, rflags.DryRun)

			case "publishImage":
				val, _, err2 := d.HashPublishImage(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				err = maybeSync(matchCtx, val, rflags.DryRun)

			case "service":
				val, _, err2 := d.HashService(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				err = maybeSync(matchCtx, val, rflags.DryRun)

			case "hostService":
				val, _, err2 := d.HashHostService(e.Value)
				if err2 != nil {
					err = err2
					break
				}
				err = maybeSync(matchCtx, val, rflags.DryRun)
			}

			if err != nil {
				if errors.Is(err, context.Canceled) || isDaggerQueryError(err) {
					fmt.Printf("ABORT: %s (%s)\n", name, kind)
				} else {
					fmt.Printf("ERROR: %s (%s) %v\n", name, kind, err)
				}
				return err
			}

			fmt.Printf(" DONE: %s (%s) (%v)\n", name, kind, time.Since(start).Round(time.Millisecond))
			return nil
		})
	}

	err = g.Wait()

	if err != nil {
		return fmt.Errorf("error during sync: %v", err)
	}

	return err
}
