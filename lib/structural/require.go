package structural

import (
	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	"fmt"
)

// require one value by another

// RequireValue uses 'require' to require a subvalue from 'from'
// by checking if values unify
func RequireValue(require, from cue.Value, opts *Options) error {
	if opts == nil {
		opts = &Options{}
	}
	return requireValue(require, from, opts)
}

// this is the recursive version that also returns
// whether the value was requireed
func requireValue(require, from cue.Value, opts *Options) error {
	switch require.IncompleteKind() {
	// require anything is like noop
	case cue.TopKind:
		return nil

	// recurse on matching labels
	case cue.StructKind:
		return requireStruct(require, from, opts)

	case cue.ListKind:
		return requireList(require, from, opts)

	default:
		return requireLeaf(require, from, opts)
	}
}

func requireStruct(require, from cue.Value, opts *Options) error {
	if k := from.IncompleteKind(); k != cue.StructKind {
		e := errors.Newf(require.Pos(), "require type '%v' does not match target value type '%v'", require.IncompleteKind(), from.IncompleteKind())
		return e
	}

	iter, _ := require.Fields(defaultWalkOptions...)
	for iter.Next() {
		s := iter.Selector()
		p := cue.MakePath(s)
		f := from.LookupPath(p)
		// fmt.Println(cnt, iter.Value(), f, f.Exists())
		// check that field exists in from. Should we be checking f.Err()?
		if f.Exists() {
			err := requireValue(iter.Value(), f, opts)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("value missing required field %q", iter.Value().Path())
		}
	}

	return nil
}

func requireList(require, from cue.Value, opts *Options) error {
	if k := from.IncompleteKind(); k != cue.ListKind {
		e := errors.Newf(require.Pos(), "require type '%v' does not match target value type '%v'", require.IncompleteKind(), from.IncompleteKind())
		return e
	}

	lpt, err := getListProcType(require)
	if err != nil {
		ce := errors.Newf(require.Pos(), "%v", err)
		return ce
	}

	_ = lpt

	// how to consider different list sizes
	// if len(require) == 1, apply to all elements
	// if len(require) > 1
	//   attributes? @require(and,or,pos)
	// maybe we don't care about length if attribute is used?

	pi, _ := require.List()
	fi, _ := from.List()

	for pi.Next() && fi.Next() {
		err := requireValue(pi.Value(), fi.Value(), opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func requireLeaf(require, from cue.Value, opts *Options) error {
	// if require is concrete, so must from
	// make sure 1 does not require int
	// but we do want int to require any num
	if require.IsConcrete() {
		if from.IsConcrete() {
			r := require.Unify(from)
			if !r.Exists() || r.Err() != nil {
				return fmt.Errorf("missing required field %q", require.Path())
			}
			return nil
		} else {
			return fmt.Errorf("missing required field %q", require.Path())
		}
	} else {
		r := require.Unify(from)
		if !r.Exists() || r.Err() != nil {
			return fmt.Errorf("missing required field %q", require.Path())
		}
		return nil
	}

	return nil
}
