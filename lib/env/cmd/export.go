package cmd

import (
	"fmt"
	"slices"

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
	}
	_, kind, _ := extractMeta(e)
	// only publish containers right now
	if slices.Contains(accepting, kind) {
		return true
	}
	return false
}

func Export(args []string, rflags flags.RootPflagpole, eflags flags.EnvPflagpole, cflags flags.Env__ExportFlagpole) error {
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

	// setup dagger & cue->dagger engine
	err = R.DaggerInit()
	d, _ := dag.NewClient(R.Ctx, R.DagClient)

	buildCtx, buildSpan := dagger.Tracer().Start(R.Ctx, "hof env export")
	defer buildSpan.End()

	fmt.Println("exporting:")
	g, groupCtx := errgroup.WithContext(buildCtx)
	if eflags.Parallel > 0 {
		g.SetLimit(eflags.Parallel)
	}

	for i, e := range matches {
		i, e := i, e
		g.Go(func() error {
			name, kind, _ := extractMeta(e)
			matchCtx, matchSpan := dagger.Tracer().Start(groupCtx, fmt.Sprintf("exporting[%d]: %s (%s)", i, name, kind))
			defer matchSpan.End()

			fmt.Printf(" - %s (%s)\n", name, kind)

			switch kind {
			case "container", "dockerBuild":
				i, err := d.Container(e.Value, eflags.NoCache)
				if err != nil {
					return err
				}
				if len(cflags.Tag) == 0 {
					err = i.ExportImage(matchCtx, fmt.Sprintf("%s:%s", name, "local"))
					if err != nil {
						return err
					}
				} else {
					for _, t := range cflags.Tag {
						// this is the "tag" annotation
						i = i.WithAnnotation("org.opencontainers.image.version", t)
						err = i.ExportImage(matchCtx, fmt.Sprintf("%s:%s", name, t))
						if err != nil {
							return err
						}
					}
				}

			case "dir":
				dir, p, err := d.Dir(e.Value, eflags.NoCache)
				if err != nil {
					return err
				}
				if len(cflags.Tag) == 0 {
					_, err = dir.Export(matchCtx, p)
					if err != nil {
						return err
					}
				} else {
					for _, t := range cflags.Tag {
						_, err = dir.Export(matchCtx, fmt.Sprintf("%s-%s", p, t))
						if err != nil {
							return err
						}
					}
				}

			case "file":
				file, p, err := d.File(e.Value, eflags.NoCache)
				if err != nil {
					return err
				}

				if len(cflags.Tag) == 0 {
					_, err = file.Export(matchCtx, p)
					if err != nil {
						return err
					}
				} else {
					for _, t := range cflags.Tag {
						_, err = file.Export(matchCtx, fmt.Sprintf("%s-%s", p, t))
						if err != nil {
							return err
						}
					}
				}

			case "exportDir":
				dir, cfg, err := d.HashExportDir(e.Value)
				if err != nil {
					return err
				}
				if len(cflags.Tag) == 0 {
					_, err = dir.Export(matchCtx, cfg.Path)
					if err != nil {
						return err
					}
				} else {
					for _, t := range cflags.Tag {
						_, err = dir.Export(matchCtx, fmt.Sprintf("%s-%s", cfg.Path, t))
						if err != nil {
							return err
						}
					}
				}

			case "exportFile":
				file, cfg, err := d.HashExportFile(e.Value)
				if err != nil {
					return err
				}
				if len(cflags.Tag) == 0 {
					_, err = file.Export(matchCtx, cfg.Path)
				} else {
					for _, t := range cflags.Tag {
						_, err = file.Export(matchCtx, fmt.Sprintf("%s-%s", cfg.Path, t))
					}
				}

			case "exportImageFile":
				i, cfg, err := d.HashExportImageFile(e.Value)
				if err != nil {
					return err
				}

				if len(cflags.Tag) == 0 {
					if len(cfg.Tags) == 0 {
						j := i.WithAnnotation("org.opencontainers.image.version", "local")
						_, err = j.Export(matchCtx, cfg.Path)
						if err != nil {
							return err
						}
					}
					for _, t := range cfg.Tags {
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						_, err = j.Export(matchCtx, fmt.Sprintf("%s-%s", cfg.Path, t))
						if err != nil {
							return err
						}
					}
				} else {
					for _, t := range cflags.Tag {
						// this is the "tag" annotation
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						_, err = j.Export(matchCtx, fmt.Sprintf("%s-%s", cfg.Path, t))
						if err != nil {
							return err
						}
					}
				}

			case "exportImage":
				i, cfg, err := d.HashExportImage(e.Value)
				if err != nil {
					return err
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
							return err
						}
					}
					for _, t := range cfg.Tags {
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						err = j.ExportImage(matchCtx, fmt.Sprintf("%s:%s", url, t))
						if err != nil {
							return err
						}
					}
				} else {
					for _, t := range cflags.Tag {
						// this is the "tag" annotation
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						err = j.ExportImage(matchCtx, fmt.Sprintf("%s:%s", url, t))
						if err != nil {
							return err
						}
					}
				}

			case "publishImage":
				i, cfg, err := d.HashPublishImage(e.Value)
				if err != nil {
					return err
				}
				url := cfg.Name
				if cfg.Reg != "" {
					url = fmt.Sprintf("%s/%s", cfg.Reg, cfg.Name)
				}

				if len(cflags.Tag) == 0 {
					if len(cfg.Tags) == 0 {
						t := "local"
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						_, err = j.Publish(matchCtx, fmt.Sprintf("%s:%s", url))
						if err != nil {
							return err
						}
					}
					for _, t := range cfg.Tags {
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						_, err = j.Publish(matchCtx, fmt.Sprintf("%s:%s", url, t))
						if err != nil {
							return err
						}
					}
				} else {
					for _, t := range cflags.Tag {
						// this is the "tag" annotation
						j := i.WithAnnotation("org.opencontainers.image.version", t)
						_, err = j.Publish(matchCtx, fmt.Sprintf("%s:%s", url, t))
						if err != nil {
							return err
						}
					}
				}

			}
			return nil
		})
	}
	err = g.Wait()

	return err
}
