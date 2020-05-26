package structural

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/cuetils"
)

func RunPickFromArgs(orig, pick string, entrypoints []string) error {
	fmt.Println("lib/st.Pick", orig, pick, entrypoints)

	crt, err := cuetils.CueRuntimeFromEntrypointsAndFlags(entrypoints)
	if err != nil {
		return err
	}

	fmt.Println("loaded")
	fmt.Println("-------------")
	lSyn, err := cuetils.ValueToSyntaxString(crt.CueValue)
	if err != nil {
		return err
	}
	fmt.Println(lSyn)
	fmt.Println("-------------")

	fmt.Println("orig")
	fmt.Println("-------------")
	oPath := strings.Split(orig, ".")
	oVal := crt.CueValue.Lookup(oPath...)
	oSyn, err := cuetils.ValueToSyntaxString(oVal)
	if err != nil {
		return err
	}
	fmt.Println(oSyn)
	fmt.Println("-------------")

	fmt.Println("pick")
	fmt.Println("-------------")
	fmt.Println(pick)
	/*
	pVal, err := crt.ConvertToValue(pick)
	if err != nil {
		return err
	}
	*/
	fmt.Println("-------------")

	fmt.Println("result")
	fmt.Println("-------------")
	rSyn, err := CuePick(oSyn, pick)
	if err != nil {
		return err
	}
	/*
	rVal, err := PickValues(oVal, pVal)
	if err != nil {
		return err
	}
	rSyn, err := cuetils.ValueToSyntaxString(rVal)
	if err != nil {
		return err
	}
	*/
	fmt.Println(rSyn)
	fmt.Println("-------------")
	return nil
}

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

func Pick(orig, pick interface{}) (cue.Value, error) {

	O, err := convertToValue(orig)
	if err != nil {
		return cue.Value{}, err
	}

	P, err := convertToValue(pick)
	if err != nil {
		return cue.Value{}, err
	}

	return PickValues(O, P)
}

func PickValues(orig, last cue.Value) (cue.Value, error) {
	out := NewpvStruct()
	err := cuePick(out, orig, last)
	if err != nil {
		return cue.Value{}, err
	}
	c, err := out.ToValue()
	return *c, err
}
