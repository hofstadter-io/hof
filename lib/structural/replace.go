package structural

import (
	"cuelang.org/go/cue"
)

func ReplaceValue(repl, val cue.Value, opts *Options) (cue.Value, error) {
	if opts == nil {
		opts = &Options{}
	}
	r, _ := replaceValue(repl, val, opts)
	return r, nil
}

func replaceValue(repl, target cue.Value, opts *Options) (cue.Value, bool) {

	switch target.IncompleteKind() {
	case cue.StructKind:
		return replaceStruct(repl, target, opts)

	case cue.ListKind:
		return replaceList(repl, target, opts)

	default:
		// should already have the same label by now
		// but maybe not if target is basic and repl is not
		return repl, true
	}
}

func replaceStruct(repl, target cue.Value, opts *Options) (cue.Value, bool) {
	ctx := target.Context()

	result := newStruct(ctx)
	iter, _ := target.Fields(defaultWalkOptions...)

	cnt := 0
	for iter.Next() {
		cnt++
		s := iter.Selector()
		p := cue.MakePath(s)
		r := repl.LookupPath(p)
		// fmt.Println(cnt, iter.Value(), f, f.Exists())
		// check that field exists in from. Should we be checking f.Err()?
		if r.Exists() {
			v, ok := replaceValue(r, iter.Value(), opts)
			// fmt.Println("r:", r, ok, p)
			if ok {
				result = result.FillPath(p, v)
			}
		} else {
			// include if not in replace
			result = result.FillPath(p, iter.Value())
		}
	}

	return result, true
}

func replaceList(repl, target cue.Value, opts *Options) (cue.Value, bool) {
	ctx := target.Context()

	ri, _ := repl.List()
	ti, _ := target.List()

	result := []cue.Value{}
	for ri.Next() && ti.Next() {
		r, ok := replaceValue(ri.Value(), ti.Value(), opts)
		if ok {
			result = append(result, r)
		}
	}
	return ctx.NewList(result...), true
}
