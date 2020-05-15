package structural

import (
	"fmt"
	"reflect"

	"cuelang.org/go/cue"
)

func CuePick(sorig, spick string) (string, error) {
	out := NewpvStruct()

	vorigi, err := r.Compile("", sorig)
	if err != nil {
		return "", err
	}
	vorig := vorigi.Value()
	if vorig.Err() != nil {
		return "", vorig.Err()
	}
	vpicki, err := r.Compile("", spick)
	if err != nil {
		return "", err
	}
	vpick := vpicki.Value()
	if vpick.Err() != nil {
		return "", vpick.Err()
	}

	err = cuePick(out, vorig, vpick)
	if err != nil {
		return "", err
	}

	return out.ToString()
}

func cuePick(out *pvStruct, vorig, vpick cue.Value) error {
	// Loop over the keys in pick
	vpickStruct, err := vpick.Struct()
	if err != nil {
		return err
	}
	vpickIter := vpickStruct.Fields()
	for vpickIter.Next() {
		key := vpickIter.Label()
		pickVal := vpickIter.Value()
		origLookup, err := vorig.LookupField(key)
		// Ignore anythig not in orig
		if err != nil {
			continue
		}
		origVal := origLookup.Value

		// If orig is a builtin and unifies with pick, then use it
		if isBuiltin(origVal) && origVal.Unify(pickVal).Kind() != cue.BottomKind {
			out.Set(key, *ExprFromValue(origVal))
			continue
		}
		if isList(origVal) {
			lval := NewpvList()
			origListIter, err := origVal.List()
			if err != nil {
				return err
			}
			// If orig is a list but pick isn't, keep all elements
			// of the list that unify with pick
			if !isList(pickVal) {
				for origListIter.Next() {
					elem := origListIter.Value()
					if elem.Unify(pickVal).Kind() != cue.BottomKind {
						lval.Append(*ExprFromValue(elem))
					}
				}
			} else if isList(pickVal) {
				// Else, consider element-wise
				pickListIter, err := pickVal.List()
				if err != nil {
					return err
				}
				for origListIter.Next() && pickListIter.Next() {
					origElem := origListIter.Value()
					pickElem := pickListIter.Value()
					if origElem.Unify(pickElem).Kind() != cue.BottomKind {
						lval.Append(*ExprFromValue(origElem))
					}
				}
			}
			out.Set(key, *lval.ToExpr())
		}

		// If orig is a struct then recurse
		if isStruct(origVal) {
			rval := NewpvStruct()
			err = cuePick(rval, origVal, pickVal)
			if err != nil {
				return err
			}
			out.Set(key, *rval.ToExpr())
		}
	}

	return nil
}

///////////

func Pick(orig, pick interface{}) (interface{}, error) {

	O, ook := orig.(cue.Value)
	P, pok := pick.(cue.Value)

	if ook && pok {
		return PickCue(O, P)
	}

	if !(ook || pok) {
		return PickGo(orig, pick)
	}

	return nil, fmt.Errorf("structural.Pick - Incompatible types %v and %v", reflect.TypeOf(orig), reflect.TypeOf(pick))
}

func PickCue(orig, pick cue.Value) (cue.Value, error) {
	fmt.Println("PickCue - no implemented")
	return cue.Value{}, nil
}

func PickGo(orig, pick interface{}) (interface{}, error) {
	fmt.Println("PickGo - no implemented")
	return nil, nil
}
