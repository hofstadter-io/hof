package structural

import (
	"fmt"
	"reflect"

	"cuelang.org/go/cue"
)

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

func Merge(orig, last interface{}) (interface{}, error) {

	O, ook := orig.(cue.Value)
	L, lok := last.(cue.Value)

	if ook && lok {
		return MergeCue(O, L)
	}

	if !(ook || lok) {
		return MergeGo(orig, last)
	}

	return nil, fmt.Errorf("structural.Merge - Incompatible types %v and %v", reflect.TypeOf(orig), reflect.TypeOf(last))
}

func MergeCue(orig, last cue.Value) (cue.Value, error) {
	fmt.Println("MergeCue - no implemented")
	return cue.Value{}, nil
}

func MergeGo(orig, last interface{}) (interface{}, error) {
	fmt.Println("MergeGo - no implemented")
	return nil, nil
}
