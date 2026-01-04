package cmd

import (
	"context"
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

	// "exportCuefig",
	// "exportDagger",
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
	d, _ := dag.NewClient(R.Ctx, R.DagClient)

	buildCtx, buildSpan := dagger.Tracer().Start(R.Ctx, "hof env sync")
	defer buildSpan.End()

	// A helper function to handle the common DryRun + Sync logic

	// do actual work
	fmt.Printf("  %v\n", time.Since(veryStart).Round(time.Millisecond))
	fmt.Println("sync'n")

	g, groupCtx := errgroup.WithContext(buildCtx)
	if eflags.Parallel > 0 {
		g.SetLimit(eflags.Parallel)
	}

	for i, e := range matches {
		i, e := i, e
		g.Go(func() error {
			name, kind, _ := extractMeta(e)
			matchCtx, matchSpan := dagger.Tracer().Start(groupCtx, fmt.Sprintf("building[%d]: %s (%s)", i, name, kind))
			defer matchSpan.End()
			fmt.Printf("%3d. %-27s (%s)", i, name, kind)
			start := time.Now()
			defer func() {
				fmt.Printf("  %v\n", time.Since(start).Round(time.Millisecond))
			}()

			// it would be freaking sweet if there was a way to get generics or something to move these syncs out and have just one, but different types and returns

			switch kind {
			case "container", "dockerBuild", "hostImage":
				val, err := d.Container(e.Value, eflags.NoCache)
				if err != nil {
					return err
				}
				return maybeSync(matchCtx, val, rflags.DryRun)

			case "dir", "gitRepo", "hostDir":
				val, _, err := d.Dir(e.Value, eflags.NoCache)
				if err != nil {
					return err
				}
				return maybeSync(matchCtx, val, rflags.DryRun)

			case "file", "hostFile":
				val, _, err := d.File(e.Value, eflags.NoCache)
				if err != nil {
					return err
				}
				return maybeSync(matchCtx, val, rflags.DryRun)

			case "exportDir":
				val, _, err := d.HashExportDir(e.Value)
				if err != nil {
					return err
				}
				if rflags.DryRun {
					return nil
				}
				_, err = val.Sync(matchCtx)
				return err

			case "exportFile":
				val, _, err := d.HashExportFile(e.Value)
				if err != nil {
					return err
				}
				return maybeSync(matchCtx, val, rflags.DryRun)

			case "exportImageFile":
				val, _, err := d.HashExportImageFile(e.Value)
				if err != nil {
					return err
				}
				return maybeSync(matchCtx, val, rflags.DryRun)

			case "exportImage":
				val, _, err := d.HashExportImage(e.Value)
				if err != nil {
					return err
				}
				return maybeSync(matchCtx, val, rflags.DryRun)

			case "publishImage":
				val, _, err := d.HashPublishImage(e.Value)
				if err != nil {
					return err
				}
				return maybeSync(matchCtx, val, rflags.DryRun)

			case "service":
				val, _, err := d.HashService(e.Value)
				if err != nil {
					return err
				}
				return maybeSync(matchCtx, val, rflags.DryRun)

			case "hostService":
				val, _, err := d.HashHostService(e.Value)
				if err != nil {
					return err
				}
				return maybeSync(matchCtx, val, rflags.DryRun)
			}

			return nil
		})
	}
	err = g.Wait()

	if err != nil {
		return fmt.Errorf("error during sync: %v", err)
	}

	return err
}
