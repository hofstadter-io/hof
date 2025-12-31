package cmd

import (
	"fmt"
	"slices"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/env"
	"github.com/hofstadter-io/hof/lib/env/dag"
)

func exportable(e *env.Env) bool {
	accepting := []string{
		"container", "hostImage", "dockerBuild",
		"dir", "hostDir", "gitRepo",
		"file", "hostFile",
		"exportFile", "exportDir",
		"exportImage", "exportImageFile", "publishImage",
	}
	_, kind := extractMeta(e)
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

	fmt.Println("exporting:")
	for _, e := range matches {
		name, kind := extractMeta(e)
		fmt.Printf(" - %s (%s)\n", name, kind)

		switch kind {
		case "container", "hostImage", "dockerBuild":
			i, err := d.Container(e, eflags.NoCache)
			if err != nil {
				return err
			}
			if len(cflags.Tag) == 0 {
				err = i.ExportImage(R.Ctx, fmt.Sprintf("%s:%s", name, "local"))
				if err != nil {
					return err
				}
			} else {
				for _, t := range cflags.Tag {
					// this is the "tag" annotation
					i = i.WithAnnotation("org.opencontainers.image.version", t)
					err = i.ExportImage(R.Ctx, fmt.Sprintf("%s:%s", name, t))
					if err != nil {
						return err
					}
				}
			}

		case "dir", "hostDir", "gitRepo":
			dir, p, err := d.Dir(e, eflags.NoCache)
			if err != nil {
				return err
			}
			if len(cflags.Tag) == 0 {
				_, err = dir.Export(R.Ctx, p)
				if err != nil {
					return err
				}
			} else {
				for _, t := range cflags.Tag {
					_, err = dir.Export(R.Ctx, fmt.Sprintf("%s-%s", p, t))
					if err != nil {
						return err
					}
				}
			}

		case "file", "hostFile":
			file, p, err := d.File(e, eflags.NoCache)
			if err != nil {
				return err
			}

			if len(cflags.Tag) == 0 {
				_, err = file.Export(R.Ctx, p)
				if err != nil {
					return err
				}
			} else {
				for _, t := range cflags.Tag {
					_, err = file.Export(R.Ctx, fmt.Sprintf("%s-%s", p, t))
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
				_, err = dir.Export(R.Ctx, cfg.Path)
				if err != nil {
					return err
				}
			} else {
				for _, t := range cflags.Tag {
					_, err = dir.Export(R.Ctx, fmt.Sprintf("%s-%s", cfg.Path, t))
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
				_, err = file.Export(R.Ctx, cfg.Path)
			} else {
				for _, t := range cflags.Tag {
					_, err = file.Export(R.Ctx, fmt.Sprintf("%s-%s", cfg.Path, t))
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
					_, err = j.Export(R.Ctx, cfg.Path)
					if err != nil {
						return err
					}
				}
				for _, t := range cfg.Tags {
					j := i.WithAnnotation("org.opencontainers.image.version", t)
					_, err = j.Export(R.Ctx, fmt.Sprintf("%s-%s", cfg.Path, t))
					if err != nil {
						return err
					}
				}
			} else {
				for _, t := range cflags.Tag {
					// this is the "tag" annotation
					j := i.WithAnnotation("org.opencontainers.image.version", t)
					_, err = j.Export(R.Ctx, fmt.Sprintf("%s-%s", cfg.Path, t))
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
					err = j.ExportImage(R.Ctx, fmt.Sprintf("%s:%s", url, t))
					if err != nil {
						return err
					}
				}
				for _, t := range cfg.Tags {
					j := i.WithAnnotation("org.opencontainers.image.version", t)
					err = j.ExportImage(R.Ctx, fmt.Sprintf("%s:%s", url, t))
					if err != nil {
						return err
					}
				}
			} else {
				for _, t := range cflags.Tag {
					// this is the "tag" annotation
					j := i.WithAnnotation("org.opencontainers.image.version", t)
					err = j.ExportImage(R.Ctx, fmt.Sprintf("%s:%s", url, t))
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
					_, err = j.Publish(R.Ctx, fmt.Sprintf("%s:%s", url))
					if err != nil {
						return err
					}
				}
				for _, t := range cfg.Tags {
					j := i.WithAnnotation("org.opencontainers.image.version", t)
					_, err = j.Publish(R.Ctx, fmt.Sprintf("%s:%s", url, t))
					if err != nil {
						return err
					}
				}
			} else {
				for _, t := range cflags.Tag {
					// this is the "tag" annotation
					j := i.WithAnnotation("org.opencontainers.image.version", t)
					_, err = j.Publish(R.Ctx, fmt.Sprintf("%s:%s", url, t))
					if err != nil {
						return err
					}
				}
			}

		}

	}

	return err
}
