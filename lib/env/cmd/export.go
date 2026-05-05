package cmd

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"dagger.io/dagger"
	"golang.org/x/sync/errgroup"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/env/dag"
)

func exportable(e *env.Env) bool {
	accepting := []string{
		"container", "dockerBuild",
		"dir", "file",
		"exportFile", "exportDir",
		"exportImage", "exportImageFile", "publishImage",
		"hostExec",
	}
	_, kind, _ := extractMeta(e)
	// only publish containers right now
	if slices.Contains(accepting, kind) {
		return true
	}
	return false
}

func Export(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole, cflags flags.Env__ExportFlagpole) error {
	veryStart := time.Now()
	// some quick setup and early filtering
	R, matches, err := commonStart(args, rflags, eflags, exportable)
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
	d, _ := dag.NewClient(R.Ctx, R.DAG)

	buildCtx, buildSpan := dagger.Tracer().Start(R.Ctx, "hof env export")
	defer buildSpan.End()

	fmt.Printf("  %v\n", time.Since(veryStart).Round(time.Millisecond))

	fmt.Println("exporting:")
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

	for ii, ee := range matches {
		i, e := ii, ee
		g.Go(func() error {
			name, kind, _ := extractMeta(e)
			matchCtx, matchSpan := dagger.Tracer().Start(groupCtx, fmt.Sprintf("exporting[%d]: %s (%s)", i, name, kind))
			defer matchSpan.End()

			if groupCtx.Err() != nil {
				fmt.Printf("ABORT: %s (%s)\n", name, kind)
				return groupCtx.Err()
			}

			fmt.Printf("START: %s (%s)\n", name, kind)
			start := time.Now()

			var err error
			switch kind {
			case "container", "dockerBuild":
				i, err2 := d.Container(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				if len(cflags.Tag) == 0 {
					err = i.ExportImage(matchCtx, fmt.Sprintf("%s:%s", name, "local"))
					if err != nil {
						break
					}
				} else {
					for _, t := range cflags.Tag {
						// this is the "tag" annotation
						i = i.WithAnnotation("org.opencontainers.image.version", t)
						err = i.ExportImage(matchCtx, fmt.Sprintf("%s:%s", name, t))
						if err != nil {
							break
						}
					}
				}

			case "dir":
				dir, p, err2 := d.Dir(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				if len(cflags.Tag) == 0 {
					_, err = dir.Export(matchCtx, p)
					if err != nil {
						break
					}
				} else {
					for _, t := range cflags.Tag {
						_, err = dir.Export(matchCtx, fmt.Sprintf("%s-%s", p, t))
						if err != nil {
							break
						}
					}
				}

			case "file":
				file, p, err2 := d.File(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}

				if len(cflags.Tag) == 0 {
					_, err = file.Export(matchCtx, p)
					if err != nil {
						break
					}
				} else {
					for _, t := range cflags.Tag {
						_, err = file.Export(matchCtx, fmt.Sprintf("%s-%s", p, t))
						if err != nil {
							break
						}
					}
				}

			case "exportDir":
				dir, cfg, err2 := d.HashExportDir(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				if len(cflags.Tag) == 0 {
					_, err = dir.Export(matchCtx, cfg.Path)
					if err != nil {
						break
					}
				} else {
					for _, t := range cflags.Tag {
						_, err = dir.Export(matchCtx, fmt.Sprintf("%s-%s", cfg.Path, t), dagger.DirectoryExportOpts{
							Wipe: cfg.Wipe,
						})
						if err != nil {
							break
						}
					}
				}

			case "exportFile":
				file, cfg, err2 := d.HashExportFile(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				if len(cflags.Tag) == 0 {
					_, err = file.Export(matchCtx, cfg.Path)
				} else {
					for _, t := range cflags.Tag {
						_, err = file.Export(matchCtx, fmt.Sprintf("%s-%s", cfg.Path, t))
						if err != nil {
							break
						}
					}
				}

			case "exportImageFile":
				i, cfg, err2 := d.HashExportImageFile(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}

				if len(cflags.Tag) == 0 {
					if len(cfg.Tags) == 0 {
						j := i.WithAnnotation("org.opencontainers.image.version", "local")
						_, err = j.Export(matchCtx, cfg.Path)
						if err != nil {
							break
						}
					}
					for _, t := range cfg.Tags {
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						_, err = j.Export(matchCtx, fmt.Sprintf("%s-%s", cfg.Path, t))
						if err != nil {
							break
						}
					}
				} else {
					for _, t := range cflags.Tag {
						// this is the "tag" annotation
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						_, err = j.Export(matchCtx, fmt.Sprintf("%s-%s", cfg.Path, t))
						if err != nil {
							break
						}
					}
				}

			case "exportImage":
				i, cfg, err2 := d.HashExportImage(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				url := cfg.Name
				if cfg.Reg != "" {
					url = fmt.Sprintf("%s/%s", cfg.Reg, cfg.Name)
				}

				if len(cflags.Tag) == 0 {
					if len(cfg.Tags) == 0 {
						t := "local"
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						err = j.ExportImage(matchCtx, fmt.Sprintf("%s:%s", url, t))
						if err != nil {
							break
						}
					}
					for _, t := range cfg.Tags {
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						err = j.ExportImage(matchCtx, fmt.Sprintf("%s:%s", url, t))
						if err != nil {
							break
						}
					}
				} else {
					for _, t := range cflags.Tag {
						// this is the "tag" annotation
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						err = j.ExportImage(matchCtx, fmt.Sprintf("%s:%s", url, t))
						if err != nil {
							break
						}
					}
				}

			case "publishImage":
				i, cfg, err2 := d.HashPublishImage(e.Value, eflags.NoCache)
				if err2 != nil {
					err = err2
					break
				}
				url := cfg.Name
				if cfg.Reg != "" {
					url = fmt.Sprintf("%s/%s", cfg.Reg, cfg.Name)
				}

				if len(cflags.Tag) == 0 {
					if len(cfg.Tags) == 0 {
						t := "local"
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						_, err = j.Publish(matchCtx, fmt.Sprintf("%s:%s", url, "latest"))
						if err != nil {
							break
						}
					}
					for _, t := range cfg.Tags {
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						_, err = j.Publish(matchCtx, fmt.Sprintf("%s:%s", url, t))
						if err != nil {
							break
						}
					}
				} else {
					for _, t := range cflags.Tag {
						// this is the "tag" annotation
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						_, err = j.Publish(matchCtx, fmt.Sprintf("%s:%s", url, t))
						if err != nil {
							break
						}
					}
				}

			case "hostExec":
				err = d.HashHostExec(e.Value)

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

	return err
}
