package datamodel

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"cuelang.org/go/cue"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func RunCheckpointFromArgs(args []string, flgs flags.DatamodelPflagpole) error {
	// fmt.Println("lib/datamodel.Checkpoint", args)

	dms, err := LoadDatamodels(args, flgs)
	if err != nil {
		return err
	}

	tag := time.Now().UTC().Format(tagFmt)
	fmt.Println("creating checkpoint:", tag)

	had := false
	for _, dm := range dms {
		if dm.status == "dirty" {
			had = true
			err = checkpointDatamodel(dm, tag)
			if err != nil {
				return err
			}
		}
	}
	if !had {
		return fmt.Errorf("no datamodels needed checkpointing")
	}

	return nil
}

func checkpointDatamodel(dm *Datamodel, tag string) error {
	fmt.Println("checkpointing:", dm.Name)

	dir, err := FindHistoryBaseDir()
	if err != nil {
		return err
	}

	hdir := filepath.Join(dir, dm.Name)
	fmt.Println("hdir:", hdir)

	err = yagu.Mkdir(hdir)
	if err != nil {
		return err
	}

	str, err := cuetils.ValueToSyntaxString(
		dm.value,
		cue.Attributes(true),
		cue.Concrete(false),
		cue.Definitions(true),
		cue.Docs(true),
		cue.Hidden(true),
		cue.Final(),
		cue.Optional(false),
		cue.ResolveReferences(false),
	)

	if err != nil {
		return err
	}

	out := fmt.Sprintf("ver_%s: %s", tag, str)
	fn := filepath.Join(hdir, fmt.Sprintf("%s.cue", tag))
	err = os.WriteFile(fn, []byte(out), 0666)
	if err != nil {
		return err
	}

	return nil
}
