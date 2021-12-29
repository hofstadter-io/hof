package datamodel

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"cuelang.org/go/cue"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
	"github.com/mattn/go-zglob"
)

// YYYYMMDDHHMMSS in Golang
const tagFmt = "20060102150405"

func LoadDatamodels(entrypoints []string, flgs flags.DatamodelPflagpole) (dms []*Datamodel, err error) {
	// load the current and all of history
	dms, err = loadDatamodelsAt(entrypoints, flgs)
	if err != nil {
		return nil, err
	}

	return dms, nil
}

func filterDatamodelsByVersion(dms []*Datamodel, flgs flags.DatamodelPflagpole) ([]*Datamodel, error) {

	// filter history
	for i, dm := range dms {
		keep := []*Datamodel{}
		for _, p := range dm.History.Past {
			// filter newer
			if flgs.Until != "" && p.version >= flgs.Until {
				// set current if it matches until
				if p.version == flgs.Until {
					p.History = &History{
						Curr: p,
					}
					dms[i] = p
				}
				continue
			}
			// filter older
			if flgs.Since != "" && p.version < flgs.Since {
				continue
			}
			keep = append(keep, p)
		}

		// set filtered history
		dms[i].History.Past = keep

	}

	return dms, nil
}

func loadDatamodelsAt(entrypoints []string, flgs flags.DatamodelPflagpole) ([]*Datamodel, error) {
	dms := []*Datamodel{}

	tag := time.Now().UTC().Format(tagFmt)

	crt, err := cuetils.CueRuntimeFromEntrypointsAndFlags(entrypoints)
	if err != nil {
		return dms, err
	}

	kvs, err := cuetils.GetByAttrKeys(crt.CueValue, "datamodel", nil, nil)
	if err != nil {
		return dms, err
	}

	for _, kv := range kvs {
		// decode and setup datamodel
		var dm Datamodel
		err = kv.Val.Decode(&dm)
		if err != nil {
			return dms, err
		}
		dm.label = kv.Key
		dm.version = "dirty-" + tag // set to current timestamp

		// make sure current value is processed same as checkpointed
		str, err := cuetils.ValueToSyntaxString(
			kv.Val,
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
			return dms, err
		}

		dm.value = crt.CueContext.CompileString(str)

		// go deeper to extract model values
		ms := dm.value.LookupPath(cue.ParsePath("Models"))
		if ms.Err() != nil {
			return dms, ms.Err()
		}
		iter, err := ms.Fields()
		if err != nil {
			return dms, err
		}

		i := 0
		for iter.Next() {
			label := iter.Selector().String()
			m, ok := dm.Models[label]
			if !ok {
				panic("cannot find label in models")
			}
			dm.Ordered = append(dm.Ordered, m)
			m.label = iter.Selector().String()
			m.value = iter.Value()
			i++
		}

		// add dm to result
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
		err = loadDatamodelHistory(dm, crt)
		if err != nil {
			return dms, nil
		}

		if len(dm.History.Past) == 0 {
			dm.status = "dirty"
		}
	}

	return dms, nil
}

func loadDatamodelHistory(dm *Datamodel, crt *cuetils.CueRuntime) error {

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
	crt.Entrypoints = entrypoints
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
		d.label = label
		d.value = value
		d.status = "ok"

		// go deeper to extract model values
		ms := d.value.LookupPath(cue.ParsePath("Models"))
		if ms.Err() != nil {
			return ms.Err()
		}
		iter, err := ms.Fields()
		if err != nil {
			return err
		}

		i := 0
		for iter.Next() {
			label := iter.Selector().String()
			m, ok := d.Models[label]
			if !ok {
				panic("cannot find label in models")
			}
			d.Ordered = append(dm.Ordered, m)
			m.label = iter.Selector().String()
			m.value = iter.Value()
			i++
		}

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
