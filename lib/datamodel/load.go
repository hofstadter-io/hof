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
	"github.com/hofstadter-io/hof/lib/structural"
	"github.com/mattn/go-zglob"
)

// YYYYMMDDHHMMSS in Golang
const tagFmt = "20060102150405"

func PrepDatamodels(entrypoints []string, flgs flags.DatamodelPflagpole) (dms []*Datamodel, err error) {

	// Loadup our Cue files
	dms, err = LoadDatamodels(entrypoints, flgs)
	if err != nil {
		return dms, err
	}

	dms, err = filterDatamodelsByTimestamp(dms, flgs)
	if err != nil {
		return dms, err
	}

	for _, dm := range dms {
		if len(dm.History.Past) == 0 {
			dm.Status = "no history"
		} else {
			past := dm.History.Past[0]
			if flgs.Since != "" {
				past = dm.History.Past[len(dm.History.Past)-1]
			}
			dm.History.Prev = past

			diff, err := structural.DiffValue(past.Value, dm.Value, nil)
			if err != nil {
				return dms, err
			}
			dm.Diff = diff

			if !diff.Exists() {
				dm.Status = "ok"
			} else {
				dm.Status = "dirty"
			}
		}
	}

	return dms, nil
}

func LoadDatamodels(entrypoints []string, flgs flags.DatamodelPflagpole) (dms []*Datamodel, err error) {
	// load the current and all of history
	dms, err = loadDatamodelsAt(entrypoints, flgs)
	if err != nil {
		return nil, err
	}

	return dms, nil
}

func filterDatamodelsByTimestamp(dms []*Datamodel, flgs flags.DatamodelPflagpole) ([]*Datamodel, error) {

	// filter history
	for i, dm := range dms {
		keep := []*Datamodel{}
		for _, p := range dm.History.Past {
			// filter newer
			if flgs.Until != "" && p.Timestamp >= flgs.Until {
				// set current if it matches until
				if p.Timestamp == flgs.Until {
					p.History = &History{
						Curr: p,
					}
					dms[i] = p
				}
				continue
			}
			// filter older
			if flgs.Since != "" && p.Timestamp < flgs.Since {
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

	kvs, err := cuetils.GetByAttrKeys(crt.CueValue, "hof", []string{"datamodel"}, nil)
	if err != nil {
		return dms, cuetils.ExpandCueError(err)
	}

	for _, kv := range kvs {
		// decode and setup datamodel
		var dm Datamodel
		err = kv.Val.Decode(&dm)
		if err != nil {
			return dms, cuetils.ExpandCueError(err)
		}
		dm.Label = kv.Key
		dm.Timestamp = "dirty-" + tag // set to current timestamp

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
			return dms, cuetils.ExpandCueError(err)
		}

		dm.Value = crt.CueContext.CompileString(str)

		// go deeper to extract model values
		// TODO, support lookup based on attribute
		ms := dm.Value.LookupPath(cue.ParsePath("Models"))
		if ms.Err() != nil {
			return dms, cuetils.ExpandCueError(ms.Err())
		}
		iter, err := ms.Fields()
		if err != nil {
			return dms, cuetils.ExpandCueError(err)
		}

		i := 0
		for iter.Next() {
			label := iter.Selector().String()
			m, ok := dm.Models[label]
			if !ok {
				panic("cannot find label in models")
			}
			dm.OrderedModels = append(dm.OrderedModels, m)
			m.Label = iter.Selector().String()
			m.Value = iter.Value()
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
			return dms, err
		}

		if len(dm.History.Past) == 0 {
			dm.Status = "dirty"
		} else {
			dm.Subsume = dm.History.Past[0].Value.Subsume(dm.Value)
		}
		// TODO(subsume), decend into Models and Fields for diff / subsume for more granular information
	}

	return dms, nil
}

func loadDatamodelHistory(dm *Datamodel, crt *cuetils.CueRuntime) error {
	// init history
	dm.History = &History{
		Curr: dm,
	}

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

	// iterate over fields (checkpoints)
	vers := []*Datamodel{}
	iter, err := crt.CueValue.Fields(cue.Attributes(true))
	if err != nil {
		return cuetils.ExpandCueError(err)
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
			return cuetils.ExpandCueError(err)
		}

		// set extra values
		d.Timestamp = tag
		d.Label = label
		d.Value = value
		d.Status = "ok"

		// extract version
		attrs := d.Value.Attributes(cue.ValueAttr)
		found := false
		for _, attr := range attrs {
			if attr.Name() == "dm_ver" {
				found = true
				d.Version, err = attr.String(0)
				if err != nil {
					return cuetils.ExpandCueError(err)
				}
			}
		}
		if !found {
			return fmt.Errorf("missing '@dm_ver' in %s @ %s", dm.Name, d.Timestamp)
		}

		//
		// TODO(subsume), decend into Models and Fields for diff / subsume for more granular information
		//

		// go deeper to extract model values
		ms := d.Value.LookupPath(cue.ParsePath("Models"))
		if ms.Err() != nil {
			return cuetils.ExpandCueError(ms.Err())
		}
		iter, err := ms.Fields(cue.Attributes(true))
		if err != nil {
			return cuetils.ExpandCueError(err)
		}

		i := 0
		for iter.Next() {
			label := iter.Selector().String()
			m, ok := d.Models[label]
			if !ok {
				panic("cannot find label in models")
			}
			d.OrderedModels = append(dm.OrderedModels, m)
			m.Label = iter.Selector().String()
			m.Value = iter.Value()
			i++
		}

		// add to history
		vers = append(vers, &d)
	}

	// sort history reverse chron
	sort.Slice(vers, func(i, j int) bool {
		return vers[i].Timestamp > vers[j].Timestamp
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
