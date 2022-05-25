package structural

import (
	"cuelang.org/go/cue"
	// "cuelang.org/go/cue/errors"
)

func MaskValue(mask, val cue.Value, opts *Options) (cue.Value, error) {
	if opts == nil {
		opts = &Options{}
	}
	r, _ := maskValue(mask, val, opts)
	return r, nil
}

// returns a value and if it should be kept / part of the return
func maskValue(mask, from cue.Value, opts *Options) (cue.Value, bool) {
	switch mask.IncompleteKind() {
	// mask everything
	case cue.TopKind:
		return newStruct(from.Context()), false

	// recurse on matching labels
	case cue.StructKind:
		return maskStruct(mask, from, opts)

	case cue.ListKind:
		return maskList(mask, from, opts)
	// (basic lit types)
	default:
		return maskLeaf(mask, from, opts)
	}

}

func maskStruct(mask, from cue.Value, opts *Options) (cue.Value, bool) {
	ctx := mask.Context()
	// eflags := opts.AllTypeErrors || opts.NodeTypeErrors

	if mask.IncompleteKind() != from.IncompleteKind() {
		return from, true
	}
	//if eflags && mask.IncompleteKind() != from.IncompleteKind() {
	//// emsg, hasErrMsg := getErrorAttrMsg(mask)
	//e := errors.Newf(mask.Pos(), "mask type '%v' does not match target value type '%v'", mask.IncompleteKind(), from.IncompleteKind())
	//ev := ctx.MakeError(e)
	//return ev, true
	//}

	result := newStruct(ctx)
	iter, _ := from.Fields(defaultWalkOptions...)

	cnt := 0
	for iter.Next() {
		cnt++
		s := iter.Selector()
		p := cue.MakePath(s)
		m := mask.LookupPath(p)
		// fmt.Println(cnt, iter.Value(), f, f.Exists())
		// check that field exists in from. Should we be checking f.Err()?
		if m.Exists() {
			r, include := maskValue(m, iter.Value(), opts)
			// fmt.Println("r:", r, ok, p)
			if include {
				result = result.FillPath(p, r)
			}
		} else {
			// include if not in mask
			result = result.FillPath(p, iter.Value())
		}
	}

	// need to check for {...}
	// no fields and open
	if cnt == 0 && mask.Allows(cue.AnyString) {
		return from, true
	}

	// fmt.Println("result:", result)

	return result, true

}

func maskList(mask, from cue.Value, opts *Options) (cue.Value, bool) {
	ctx := mask.Context()

	if k := from.IncompleteKind(); k != cue.ListKind {
		// should this return or just continue? do we need some way to specify?
		// probably prefer to be more strict, so that you know your schemas
		// return errors.Newf(from.Pos(), "expected list, but got %v", k), true
		return from, true
	}

	//lpt, err := getListProcType(mask)
	//if err != nil {
	//ce := errors.Newf(mask.Pos(), "%v", err)
	//ev := ctx.MakeError(ce)
	//return ev, true
	//}

	// _ = lpt

	mi, _ := mask.List()
	fi, _ := from.List()

	result := []cue.Value{}
	for mi.Next() && fi.Next() {
		p, include := maskValue(mi.Value(), fi.Value(), opts)
		if include {
			result = append(result, p)
		}
	}

	return ctx.NewList(result...), true
}

// returns a value and whether it should be included
func maskLeaf(mask, from cue.Value, opts *Options) (cue.Value, bool) {
	if mask.IsConcrete() {
		if from.IsConcrete() {
			r := mask.Unify(from)
			// need to check for errors here?
			return r, !r.Exists()
		} else {
			return cue.Value{}, false
		}
	} else {
		r := mask.Unify(from)
		// need to check for errors here?
		return r, !r.Exists()
	}

}
