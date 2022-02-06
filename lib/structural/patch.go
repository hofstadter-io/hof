package structural

import (
	"cuelang.org/go/cue"
)

func PatchValue(patch, val cue.Value, opts *Options) (cue.Value, error) {
	if opts == nil {
		opts = &Options{}
	}
	r, _ := patchValue(patch, val, opts)
	return r, nil
}

func patchValue(patch, val cue.Value, opts *Options) (cue.Value, bool) {
	switch val.IncompleteKind() {
	case cue.StructKind:
		// fmt.Println("struct", orig, next)
		return patchStruct(patch, val, opts)

	case cue.ListKind:
		// fmt.Println("list", orig, next)
		return patchList(patch, val, opts)

	default:
		panic("should not get here")
	}
}

func patchStruct(patch, val cue.Value, opts *Options) (cue.Value, bool) {
	ctx := val.Context()
	result := newStruct(ctx)
	rmv := patch.LookupPath(cue.ParsePath("\"-\""))

	iter, _ := val.Fields(defaultWalkOptions...)
	for iter.Next() {
		s := iter.Selector()
		p := cue.MakePath(s)
		r := rmv.LookupPath(p)

		// skip any field which exists, it was the same in the original diff
		// so if not in remove...
		if !r.Exists() {
			v := patch.LookupPath(p)
			if v.Exists() {
				// if in patch, we need to recurse
				r, _ := patchValue(v, iter.Value(), opts)
				result = result.FillPath(p, r)
			} else {
				// keep any field which !exists, it was the same in the original diff
				result = result.FillPath(p, iter.Value())
			}
		}
	}

	add := patch.LookupPath(cue.ParsePath("\"+\""))
	if add.Exists() {
		result = result.Unify(add)
	}

	return result, true
}

func patchList(patch, val cue.Value, opts *Options) (cue.Value, bool) {
	ctx := val.Context()
	oi, _ := patch.List()
	ni, _ := val.List()

	result := []cue.Value{}
	for oi.Next() && ni.Next() {
		v, ok := patchValue(oi.Value(), ni.Value(), opts)
		if ok {
			result = append(result, v)
		}
	}

	return ctx.NewList(result...), true
}
