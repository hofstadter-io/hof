package structural

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
)

// PickValue uses 'pick' to pick a subvalue from 'from'
// by checking if values unify
func PickValue(pick, from cue.Value, opts *Options) (cue.Value, error) {
	if opts == nil {
		opts = &Options{}
	}
	p, _ := pickValue(pick, from, opts)
	return p, nil
}

// this is the recursive version that also returns
// whether the value was picked
func pickValue(pick, from cue.Value, opts *Options) (cue.Value, bool) {
	switch pick.IncompleteKind() {
	// pick everything
	case cue.TopKind:
		return from, true

	// recurse on matching labels
	case cue.StructKind:
		return pickStruct(pick, from, opts)

	case cue.ListKind:
		return pickList(pick, from, opts)

	default:
		return pickLeaf(pick, from, opts)
	}
}

func pickStruct(pick, from cue.Value, opts *Options) (cue.Value, bool) {
	ctx := pick.Context()

	if k := from.IncompleteKind(); k != cue.StructKind {
		if opts.NodeTypeErrors {
			e := errors.Newf(pick.Pos(), "pick type '%v' does not match target value type '%v'", pick.IncompleteKind(), from.IncompleteKind())
			ev := ctx.MakeError(e)
			return ev, true
		}
	}

	result := newStruct(ctx)
	iter, _ := pick.Fields(defaultWalkOptions...)

	cnt := 0
	for iter.Next() {
		cnt++
		s := iter.Selector()
		p := cue.MakePath(s)
		f := from.LookupPath(p)
		// fmt.Println(cnt, iter.Value(), f, f.Exists())
		// check that field exists in from. Should we be checking f.Err()?
		if f.Exists() {
			r, ok := pickValue(iter.Value(), f, opts)
			// fmt.Println("r:", r, ok, p)
			if ok {
				result = result.FillPath(p, r)
			}
		}
	}

	// need to check for {...}
	// no fields and open
	if cnt == 0 && pick.Allows(cue.AnyString) {
		return from, true
	}

	// fmt.Println("result:", result)

	return result, true

	return pick, false
}

func pickList(pick, from cue.Value, opts *Options) (cue.Value, bool) {
	ctx := pick.Context()

	if k := from.IncompleteKind(); k != cue.ListKind {
		// should this return or just continue? do we need some way to specify?
		// probably prefer to be more strict, so that you know your schemas
		// return errors.Newf(from.Pos(), "expected list, but got %v", k), true
		return newStruct(ctx), false
	}

	lpt, err := getListProcType(pick)
	if err != nil {
		ce := errors.Newf(pick.Pos(), "%v", err)
		ev := ctx.MakeError(ce)
		return ev, true
	}

	_ = lpt

	// how to consider different list sizes
	// if len(pick) == 1, apply to all elements
	// if len(pick) > 1
	//   attributes? @pick(and,or,pos)
	// maybe we don't care about length if attribute is used?

	pi, _ := pick.List()
	fi, _ := from.List()

	result := []cue.Value{}
	for pi.Next() && fi.Next() {
		p, ok := pickValue(pi.Value(), fi.Value(), opts)
		if ok {
			result = append(result, p)
		}
	}

	return ctx.NewList(result...), true
}

func pickLeaf(pick, from cue.Value, opts *Options) (cue.Value, bool) {
	// if pick is concrete, so must from
	// make sure 1 does not pick int
	// but we do want int to pick any num
	if pick.IsConcrete() {
		if from.IsConcrete() {
			r := pick.Unify(from)
			return r, r.Exists()
		} else {
			return cue.Value{}, false
		}
	} else {
		r := pick.Unify(from)
		return r, r.Exists()
	}

	return pick, false
}
