package structural

import (
	"fmt"

	"cuelang.org/go/cue"
)

func MergeValues(orig, update cue.Value) (cue.Value, error) {
	out := NewpvStruct()
	err := cueMerge(out, orig, update)
	if err != nil {
		return cue.Value{}, err
	}
	c, err := out.ToValue()
	return *c, err
}

func RunMergeFromArgs(orig, update string, entrypoints []string) error {
	fmt.Println("lib/st.Merge", orig, update, entrypoints)

	return nil
}

func CueMerge(sorig, snew string) (string, error) {
	out := NewpvStruct()

	vorigi, err := r.Compile("", sorig)
	if err != nil {
		return "", err
	}
	vorig := vorigi.Value()
	if vorig.Err() != nil {
		return "", vorig.Err()
	}
	vnewi, err := r.Compile("", snew)
	if err != nil {
		return "", err
	}
	vnew := vnewi.Value()
	if vnew.Err() != nil {
		return "", vnew.Err()
	}

	err = cueMerge(out, vorig, vnew)
	if err != nil {
		return "", err
	}

	return out.ToString()
}

func cueMerge(out *pvStruct, vorig, vnew cue.Value) error {
	// Loop over keys in orig
	vorigStruct, err := vorig.Struct()
	if err != nil {
		return err
	}
	vorigIter := vorigStruct.Fields()
	for vorigIter.Next() {
		origVal := vorigIter.Value()
		// Add anything that doesn't exist in new
		newVal, err := vnew.LookupField(vorigIter.Label())
		if err != nil {
			out.Set(vorigIter.Label(), *ExprFromValue(origVal))
			continue
		}

		// If the values are not the same overall shape, then fail
		if (isBuiltin(origVal) && !isBuiltin(newVal.Value)) ||
			(isStruct(origVal) && !isStruct(newVal.Value)) ||
			(isList(origVal) && !isList(newVal.Value)) {
			return fmt.Errorf("invalid merge: %s has different type", vorigIter.Label())
		}

		// Always use new for builtins and lists
		// TODO handle lists better
		if isBuiltin(origVal) || isList(origVal) {
			out.Set(vorigIter.Label(), *ExprFromValue(newVal.Value))
			continue
		}

		// Recurse for structs
		if !isStruct(origVal) || !isStruct(newVal.Value) {
			panic("should not reach")
		}
		rval := NewpvStruct()
		err = cueMerge(rval, origVal, newVal.Value)
		if err != nil {
			return err
		}
		out.Set(vorigIter.Label(), *rval.ToExpr())
	}

	// Loop over keys in new
	vnewStruct, err := vnew.Struct()
	if err != nil {
		return err
	}
	vnewIter := vnewStruct.Fields()
	for vnewIter.Next() {
		// Add any new keys not in orig
		_, err := vorig.LookupField(vnewIter.Label())
		if err != nil {
			out.Set(vnewIter.Label(), *ExprFromValue(vnewIter.Value()))
		}
	}

	return nil

}

///////////

func Merge(orig, last interface{}) (cue.Value, error) {

	O, ook := orig.(cue.Value)
	if !ook {
		switch T := orig.(type) {
		case string:
			i, err := r.Compile("", orig)
			if err != nil {
				return O, err
			}
			v := i.Value()
			if v.Err() != nil {
				return v, v.Err()
			}
			O = v

		default:
			return O, fmt.Errorf("unknown type %v in Merge(orig,_)", T)
		}
	}

	L, lok := last.(cue.Value)
	if !lok {
		switch T := last.(type) {
		case string:
			i, err := r.Compile("", last)
			if err != nil {
				return O, err
			}
			v := i.Value()
			if v.Err() != nil {
				return v, v.Err()
			}
			L = v

		default:
			return L, fmt.Errorf("unknown type %v in Merge(_,last)", T)
		}
	}

	return MergeValues(O, L)
}
