package datamodel

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"cuelang.org/go/cue"
	"github.com/Masterminds/semver/v3"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/hofstadter-io/hof/lib/yagu"
)

func RunCheckpointFromArgs(args []string, bump string, flgs flags.DatamodelPflagpole) error {
	// fmt.Println("lib/datamodel.Checkpoint", args)

	dms, err := PrepDatamodels(args, flgs)
	if err != nil {
		return err
	}

	timestamp := time.Now().UTC().Format(tagFmt)
	fmt.Printf("creating %q checkpoint: %s\n", bump, timestamp)

	had := false
	for _, dm := range dms {
		if dm.Status == "dirty" || dm.Status == "no history" {
			had = true
			err = checkpointDatamodel(dm, timestamp, bump)
			if err != nil {
				return err
			}
			fmt.Printf(" + %s @ %s\n", dm.Name, dm.Version)
		}
	}
	if !had {
		return fmt.Errorf("no datamodels needed checkpointing")
	}

	return nil
}

func checkpointDatamodel(dm *Datamodel, timestamp, bump string) error {
	// check subsumption
	if dm.Subsume != nil {
		fmt.Println("subsume:", dm.Subsume)
		if bump == "patch" {
			return fmt.Errorf("\nBackwards incompatible changes, you should use a major or minor bump")
		}
	}

	tag := "0.0.0"

	if len(dm.History.Past) > 0 {
		// get last semver tag
		tag = dm.History.Past[0].Version
	}

	v, err := semver.StrictNewVersion(tag)
	if err != nil {
		return fmt.Errorf("error parsing last semver in history: %q\n%s", tag, err.Error())
	}

	// bump semver tag
	switch bump {
	case "major":
		b := v.IncMajor()
		v = &b
	case "minor":
		b := v.IncMinor()
		v = &b
	case "patch":
		b := v.IncPatch()
		v = &b
	default:
		n, err := semver.StrictNewVersion(bump)
		if err != nil {
			return fmt.Errorf("error parsing next semver: %q\n%s", bump, err.Error())
		}
		if !n.GreaterThan(v) {
			return fmt.Errorf("new version must be greater than last: %v < %v", n, v)
		}
		v = n
	}

	tag = v.String()
	dm.Version = tag

	fmt.Printf("%s bump ver to: %s\n", bump, dm.Version)

	// add to dm
	ctx := dm.Value.Context()
	str := fmt.Sprintf("{ @dm_ver(%s) }", tag)
	val := ctx.CompileString(str)
	dm.Value = dm.Value.Unify(val)

	return writeCheckpoint(dm, timestamp)
}

func writeCheckpoint(dm *Datamodel, timestamp string) error {
	dir, err := FindHistoryBaseDir()
	if err != nil {
		return err
	}

	hdir := filepath.Join(dir, dm.Name)

	err = yagu.Mkdir(hdir)
	if err != nil {
		return err
	}

	str, err := cuetils.ValueToSyntaxString(
		dm.Value,
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

	out := fmt.Sprintf("ver_%s: %s", timestamp, str)
	fn := filepath.Join(hdir, fmt.Sprintf("%s.cue", timestamp))
	err = os.WriteFile(fn, []byte(out), 0666)
	if err != nil {
		return err
	}

	return nil
}
