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
			} else {
				for _, t := range cflags.Tag {
					// this is the "tag" annotation
					i = i.WithAnnotation("org.opencontainers.image.version", t)
					err = i.ExportImage(R.Ctx, fmt.Sprintf("%s:%s", name, t))
				}
			}

		case "dir", "hostDir", "gitRepo":
			dir, p, err := d.Dir(e, eflags.NoCache)
			if err != nil {
				return err
			}
			if len(cflags.Tag) == 0 {
				_, err = dir.Export(R.Ctx, p)
			} else {
				for _, t := range cflags.Tag {
					_, err = dir.Export(R.Ctx, fmt.Sprintf("%s-%s", p, t))
				}
			}

		case "file", "hostFile":
			file, p, err := d.File(e, eflags.NoCache)
			if err != nil {
				return err
			}

			if len(cflags.Tag) == 0 {
				_, err = file.Export(R.Ctx, p)
			} else {
				for _, t := range cflags.Tag {
					_, err = file.Export(R.Ctx, fmt.Sprintf("%s-%s", p, t))
				}
			}

		}

	}

	return err
}
