package datamodel

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"cuelang.org/go/cue/load"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/mattn/go-zglob"
)

func LoadDatamodels(entrypoints []string, flgs flags.DatamodelPflagpole) ([]*Datamodel, error) {
	dms := []*Datamodel{}

	crt, err := cuetils.CueRuntimeFromEntrypointsAndFlags(entrypoints)
	if err != nil {
		return dms, err
	}

	kvs, err := cuetils.GetByAttrKeys(crt.CueValue, "datamodel", nil, nil)
	if err != nil {
		return dms, err
	}

	for _, kv := range kvs {
		var dm Datamodel
		err = kv.Val.Decode(&dm)
		if err != nil {
			return dms, err
		}
		dm.label = kv.Key
		dm.value = kv.Val
		dms = append(dms, &dm)
	}

	if len(dms) == 0 {
		return nil, fmt.Errorf("No datamodels found")
	}

	if len(flgs.Datamodels) > 0 {
		final := []*Datamodel{}
		for _, d := range flgs.Datamodels {
			for _, dm := range dms {
				if match, _ := regexp.MatchString(d, dm.Name); match {
					final = append(final, dm)
				}
			}
		}
		dms = final
	}

	for _, dm := range dms {
		err = LoadDatamodelHistory(dm)
		if err != nil {
			return dms, nil
		}

		if len(dm.History.Past) == 0 {
			dm.status = "dirty"
		}
	}

	return dms, nil
}

func LoadDatamodelHistory(dm *Datamodel) error {

	// find module root
	base, err := cuetils.FindModuleAbsPath()
	if err != nil {
		return err
	}

	// get entrypoints
	glob := filepath.Join(base, ".hof", "dm", dm.Name, "*.cue")
	entrypoints, err := zglob.Glob(glob)
	if err != nil {
		return err
	}

	// load datamodel history as CUE
	crt := &cuetils.CueRuntime{
		Entrypoints: entrypoints,
		CueConfig: &load.Config{
			ModuleRoot: base,
			Module:     "",
			Package:    "",
			Dir:        "",
			BuildTags:  []string{},
			Tests:      false,
			Tools:      false,
			DataFiles:  false,
			Overlay:    map[string]load.Source{},
		},
	}
	err = crt.Load()
	if err != nil {
		return err
	}

	// start history
	dm.History = &History{
		Curr: dm,
	}

	// iterate over fields (checkpoints)
	vers := []*Datamodel{}
	iter, err := crt.CueValue.Fields()
	if err != nil {
		return err
	}

	for iter.Next() {
		// meta fields
		label := iter.Selector().String()
		tag := strings.TrimPrefix(label, "ver_")
		value := iter.Value()

		// decode datamodel checkpoint
		var d Datamodel
		err := value.Decode(&d)
		if err != nil {
			return err
		}

		// set extra values
		d.version = tag
		d.value = value

		// add to history
		vers = append(vers, &d)
	}

	// sort history reverse chron
	sort.Slice(vers, func(i, j int) bool {
		return vers[i].version > vers[j].version
	})

	// set history
	dm.History.Past = vers
	return nil
}

func FindHistoryBaseDir() (string, error) {
	// try to find history
	dir, err := cuetils.FindModuleAbsPath()
	if err != nil {
		return "", err
	}

	// .hof dir is peer of cue.mod
	hdir := filepath.Join(dir, ".hof/dm")

	return hdir, nil
}
