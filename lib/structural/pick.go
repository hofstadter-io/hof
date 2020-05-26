package structural

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/token"

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

func PickValues(orig, pick cue.Value) (val cue.Value, err error) {
	node, err := cuepick_Values("", "start", orig.Eval(), pick.Eval())
	if err != nil {
		return val, err
	}

	rt := cue.Runtime{}
	I, err := rt.CompileExpr(node.(ast.Expr))
	if err != nil {
		return val, err
	}

	val = I.Value()

	return val, err
}

func cuepick_Values(indent, tag string, orig, pick cue.Value) (val ast.Node, err error) {
	oKind := orig.Kind()
	oLabel, _ := orig.Label()
	pKind := pick.Kind()
	pLabel, _ := pick.Label()

	fmt.Printf("%spick: %q %q %q %q %q  BEG\n", indent, tag, oLabel, oKind, pLabel, pKind)

	if oKind != pKind {
		return cuepick_Kinds(indent + "  ", tag + "/kinds", orig, pick)
	}

	switch pKind {
	case cue.StructKind:
		return cuepick_Structs(indent + "  ", tag + "/structs", orig, pick)

	default:
		return val, fmt.Errorf("Unknown Cue Type %q in cuepick_Values same type switch", pKind)
	}

	fmt.Printf("%spick: %s  END\n", indent, tag)
	return val, err
}

func cuepick_Structs(indent, tag string, orig, pick cue.Value) (val ast.Node, err error) {

	oKind := orig.Kind()
	oLabel, _ := orig.Label()
	pKind := pick.Kind()
	pLabel, _ := pick.Label()

	fmt.Printf("%sstructs: %q %q %q %q %q  BEG\n", indent, tag, oLabel, oKind, pLabel, pKind)

	fields := []ast.Decl{}
	oStruct, _ := orig.Struct()
	iter := oStruct.Fields()
	for iter.Next() {
		ol := iter.Label()
		ov := iter.Value()
		ok := ov.Kind()

		// try to look it up
		pv := pick.Lookup(ol)
		pk := pv.Kind()

		// print everything we decide on
		fmt.Printf("%s  field: %s %s %s %s %s\n", indent, ol, ov, ok, pv, pk)
		u := pv.Unify(ov)
		fmt.Printf("%s    unify: %s %s %s %s\n", indent, ol, pv, ov, u)
		fmt.Printf("%s    check: %s\n", indent, u.Equals(ov))

		// If we unified
		if u.Equals(ov) {
			fmt.Printf("%s    equals: %q %s\n", indent, ol, ov)
			syntax := ov.Syntax()
			field := &ast.Field {
				Label: ast.NewString(ol),
				Value: syntax.(ast.Expr),
				Token: token.COLON,
			}
			fields = append(fields, field)
		}

		if isBuiltin(pv) {
			fmt.Println(indent, "   Builtin! ", ol, pk, pv)
		}

		// is it not found?
		if pv.Kind() == cue.BottomKind {
			continue
		}

		switch pk {

		case pKind:
			fmt.Printf("%s    samekind: %q %s\n", indent, ol, pv)
			nv, err := cuepick_Values(indent + "  ", tag + "/pick", ov, pv)
			if err != nil {
				return val, err
			}
			field := &ast.Field {
				Label: ast.NewString(ol),
				Value: nv.(ast.Expr),
				Token: token.COLON,
			}
			fields = append(fields, field)
			// fmt.Printf("%s  return: %q %s %v\n", indent, ol, nv, err)

		default:
			fmt.Printf("%s    default: %q %s\n", indent, ol, pv)
			syntax := pv.Syntax()
			field := &ast.Field {
				Label: ast.NewString(ol),
				Value: syntax.(ast.Expr),
				Token: token.COLON,
			}
			fields = append(fields, field)
		}
	}

	val = ast.NewStruct()
	s := val.(*ast.StructLit)
	for _, f := range fields {
		s.Elts = append(s.Elts, f)
	}
	// fmt.Printf("%sstructs: %s  END %v %v\n", indent, tag, val, fields)
	return val, err
}

func cuepick_Lists(indent, tag string, orig, pick cue.Value) (val ast.Node, err error) {
	oKind := orig.Kind()
	oLabel, _ := orig.Label()
	pKind := pick.Kind()
	pLabel, _ := pick.Label()

	fmt.Printf("%slists: %q %q %q %q %q  BEG\n", indent, tag, oLabel, oKind, pLabel, pKind)



	fmt.Printf("%slists: %s  END\n", indent, tag)
	return val, err
}

func cuepick_Kinds(indent, tag string, orig, pick cue.Value) (val ast.Node, err error) {
	oKind := orig.Kind()
	oLabel, _ := orig.Label()
	pKind := pick.Kind()
	pLabel, _ := pick.Label()

	fmt.Printf("%skinds: %q %q %q %q %q  BEG\n", indent, tag, oLabel, oKind, pLabel, pKind)

	u := pick.Unify(orig)
	fmt.Println("Unify", orig, pick, u)
	if u.Equals(pick) {
		return orig.Syntax(), err
	}

	// TODO nested type switch here
	switch pKind {
	case cue.StructKind:
		return val, err

	default:
		return val, fmt.Errorf("Unknown Cue Type %q in cuepick_Values same type switch", pKind)
	}

	fmt.Printf("%skinds: %s  END\n", indent, tag)
	return val, err
}
