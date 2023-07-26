package structural

import (
	"fmt"

	"cuelang.org/go/cue"
)

func InsertValue(ins, val cue.Value, opts *Options) (cue.Value, error) {
	if opts == nil {
		opts = &Options{}
	}
	r, _ := insertValue(ins, val, opts)
	return r, nil
}

func insertValue(ins, val cue.Value, opts *Options) (cue.Value, bool) {
	switch val.IncompleteKind() {
	case cue.StructKind:
		return insertStruct(ins, val, opts)

	case cue.ListKind:
		return insertList(ins, val, opts)

	default:
		// should already have the same label by now
		// but maybe not if target is basic and repl is not
		return val, true
	}
}

func insertStruct(ins, val cue.Value, opts *Options) (cue.Value, bool) {

	result := val
	iter, _ := ins.Fields(defaultWalkOptions...)

	for iter.Next() {
		s := iter.Selector()
		// HACK, this works around a bug in CUE
		// p := cue.MakePath(s)
		p := cue.ParsePath(fmt.Sprint(s))
		v := val.LookupPath(p)

		// check that field exists in from. Should we be checking f.Err()?
		if v.Exists() {
			r, ok := insertValue(iter.Value(), v, opts)
			// fmt.Println("r:", r, ok, p)
			if ok {
				result = result.FillPath(p, r)
			}
		} else {
			// include if not in val
			result = result.FillPath(p, iter.Value())
		}
	}

	return result, true
}

func insertList(ins, val cue.Value, opts *Options) (cue.Value, bool) {
	ctx := val.Context()

	ii, _ := ins.List()
	vi, _ := val.List()

	result := []cue.Value{}
	for ii.Next() && vi.Next() {
		r, ok := insertValue(ii.Value(), vi.Value(), opts)
		if ok {
			result = append(result, r)
		}
	}
	return ctx.NewList(result...), true
}
