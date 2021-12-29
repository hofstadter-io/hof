package datamodel

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
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
	// try to find history
	dir, err := cuetils.FindModuleAbsPath()
	if err != nil {
		return err
	}

	hdir := filepath.Join(dir, ".hof/dm", dm.Name)

	fmt.Println("Module Root:", dir)
	fmt.Println("History Dir:", hdir)

	dm.History = &History{
		Curr: dm,
	}

	return nil
}
