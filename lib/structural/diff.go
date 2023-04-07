package structural

import (
	"cuelang.org/go/cue"
)

func DiffValue(orig, next cue.Value, opts *Options) (cue.Value, error) {
	if opts == nil {
		opts = &Options{}
	}
	r, ok := diffValue(orig, next, opts)
	if !ok {
		return cue.Value{}, nil
	}
	return r, nil
}

func diffValue(orig, next cue.Value, opts *Options) (cue.Value, bool) {

	switch orig.IncompleteKind() {
	case cue.StructKind:
		// fmt.Println("struct", orig, next)
		return diffStruct(orig, next, opts)

	case cue.ListKind:
		// fmt.Println("list", orig, next)
		return diffList(orig, next, opts)

	default:
		// fmt.Println("leaf", orig, next)
		return diffLeaf(orig, next, opts)
	}
}

func diffStruct(orig, next cue.Value, opts *Options) (cue.Value, bool) {
	ctx := orig.Context()
	result := newStruct(ctx)
	add := newStruct(ctx)
	rmv := newStruct(ctx)
	didAdd := false
	didRmv := false

	// first loop over val
	iter, _ := orig.Fields(defaultWalkOptions...)
	for iter.Next() {
		s := iter.Selector()
		p := cue.MakePath(s)
		u := next.LookupPath(p)

		// check that field exists in from. Should we be checking f.Err()?
		if u.Exists() {
			r, ok := diffValue(iter.Value(), u, opts)
			// fmt.Println("r:", r, ok, p)
			if ok {
				result = result.FillPath(p, r)
			}
		} else {
			// remove if orig not in next
			didRmv = true
			rmv = rmv.FillPath(p, iter.Value())
		}
	}

	// add anything in next that is not in orig
	iter, _ = next.Fields(defaultWalkOptions...)
	for iter.Next() {
		s := iter.Selector()
		p := cue.MakePath(s)
		v := orig.LookupPath(p)

		// check that field exists in from. Should we be checking f.Err()?
		if !v.Exists() {
			didAdd = true
			add = add.FillPath(p, iter.Value())
		}
	}

	if didRmv {
		result = result.FillPath(cue.ParsePath("\"-\""), rmv)
	}
	if didAdd {
		result = result.FillPath(cue.ParsePath("\"+\""), add)
	}

	// checks to see if nothing changed
	i := 0
	iter, _ = result.Fields()
	for iter.Next() {
		i++
	}
	if i == 0 {
		return result, false
	}

	return result, true
}

func diffList(orig, next cue.Value, opts *Options) (cue.Value, bool) {
	ctx := orig.Context()
	oi, _ := orig.List()
	ni, _ := next.List()

	result := []cue.Value{}
	for oi.Next() && ni.Next() {
		v, ok := diffValue(oi.Value(), ni.Value(), opts)
		if ok {
			result = append(result, v)
		}
	}

	return ctx.NewList(result...), len(result) != 0
}

func diffLeaf(orig, next cue.Value, opts *Options) (cue.Value, bool) {
	if orig.IncompleteKind() == next.IncompleteKind() {
		if orig.IsConcrete() == next.IsConcrete() {
			u := orig.Unify(next)
			if u.Err() == nil {
				return cue.Value{}, false
			}
		}
	}

	// need to know if this is a basic lit, so we know if we are changing a concrete value
	ctx := orig.Context()
	ret := ctx.CompileString("{}")
	ret = ret.FillPath(cue.ParsePath("\"-\""), orig)
	ret = ret.FillPath(cue.ParsePath("\"+\""), next)

	/*
	ctx := orig.Context()
	ret := newStruct(ctx)
	lbl := GetLabel(orig)

	// check if they are the same type and concreteness, and check if unify, if so, no need to add to diff
	if orig.IncompleteKind() == next.IncompleteKind() {
		if orig.IsConcrete() == next.IsConcrete() {
			u := orig.Unify(next)
			if u.Err() == nil {
				return ret, false
			}
		}
	}

	// otherwise, we have a diff to create
	rmv := newStruct(ctx)
	rmv = rmv.FillPath(cue.MakePath(lbl), orig)
	ret = ret.FillPath(cue.ParsePath("\"-\""), rmv)

	add := newStruct(ctx)
	add = add.FillPath(cue.MakePath(lbl), next)
	ret = ret.FillPath(cue.ParsePath("\"+\""), add)
	*/

	return ret, true
}
