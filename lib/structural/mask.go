package structural

import (
	"fmt"

	"cuelang.org/go/cue"
)

func RunMaskFromArgs(orig, mask string, entrypoints []string) error {
	fmt.Println("lib/st.Mask", orig, mask, entrypoints)

	return nil
}

func CueMask(sorig, smask string) (string, error) {
	out := NewpvStruct()

	vorigi, err := r.Compile("", sorig)
	if err != nil {
		return "", err
	}
	vorig := vorigi.Value()
	if vorig.Err() != nil {
		return "", vorig.Err()
	}
	vmaski, err := r.Compile("", smask)
	if err != nil {
		return "", err
	}
	vmask := vmaski.Value()
	if vmask.Err() != nil {
		return "", vmask.Err()
	}

	err = cueMask(out, vorig, vmask)
	if err != nil {
		return "", err
	}

	return out.ToString()
}

func cueMask(out *pvStruct, vorig, vmask cue.Value) error {
	// Loop over the keys in orig
	vorigStruct, err := vorig.Struct()
	if err != nil {
		return err
	}
	vorigIter := vorigStruct.Fields()
	for vorigIter.Next() {
		key := vorigIter.Label()
		origVal := vorigIter.Value()
		maskLookup, err := vmask.LookupField(key)
		// Include anythig not in mask
		if err != nil {
			out.Set(key, *ExprFromValue(origVal))
			continue
		}
		maskVal := maskLookup.Value

		// If orig is a builtin and doesn't unify with mask, then use it
		if isBuiltin(origVal) && origVal.Unify(maskVal).Kind() == cue.BottomKind {
			out.Set(key, *ExprFromValue(origVal))
			continue
		}
		if isList(origVal) {
			lval := NewpvList()
			origListIter, err := origVal.List()
			if err != nil {
				return err
			}
			// If orig is a list but mask isn't, keep all elements
			// of the list that don't unify with mask
			if !isList(maskVal) {
				for origListIter.Next() {
					elem := origListIter.Value()
					if elem.Unify(maskVal).Kind() == cue.BottomKind {
						lval.Append(*ExprFromValue(elem))
					}
				}
			} else if isList(maskVal) {
				// Else, consider element-wise
				maskListIter, err := maskVal.List()
				if err != nil {
					return err
				}
				for origListIter.Next() && maskListIter.Next() {
					origElem := origListIter.Value()
					maskElem := maskListIter.Value()
					if origElem.Unify(maskElem).Kind() == cue.BottomKind {
						lval.Append(*ExprFromValue(origElem))
					}
				}
			}
			out.Set(key, *lval.ToExpr())
		}

		// If orig is a struct then recurse
		if isStruct(origVal) {
			rval := NewpvStruct()
			err = cueMask(rval, origVal, maskVal)
			if err != nil {
				return err
			}
			out.Set(key, *rval.ToExpr())
		}
	}

	return nil
}
