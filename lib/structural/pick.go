package structural

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	// "cuelang.org/go/cue/token"

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
	// TODO< check for path?
	fmt.Println(pick)
	pVal, err := crt.ConvertToValue(pick)
	if err != nil {
		return err
	}
	fmt.Println("-------------")

	fmt.Println("result")
	fmt.Println("-------------")

	rVal, err := PickValues(oVal, pVal)
	if err != nil {
		return err
	}
	rSyn, err := cuetils.ValueToSyntaxString(rVal)
	if err != nil {
		return err
	}

	fmt.Println(rSyn)
	fmt.Println("-------------")
	return nil
}


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
	node, err := pickValues("", "start", orig, pick)
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

// top-level function which introspects the values and calls specialized functions for picking
func pickValues(indent, tag string, orig, pick cue.Value) (val ast.Node, err error) {
	//oLabel, _ := orig.Label()
	//oKind := orig.Kind()
	//pLabel, _ := pick.Label()
	pKind := pick.IncompleteKind()

	// fmt.Printf("%spick: %q %q %q %q %q  BEG\n", indent, tag, oLabel, oKind, pLabel, pKind)

	// They are not the same Kind
	//if oKind != pKind {
		//return pickKinds(indent + "  ", tag + "/kinds", orig, pick)
	//}

	// Switch on the common type
	// (list of Cue types: https://pkg.go.dev/cuelang.org/go@v0.3.0-alpha4/cue#pkg-constants)
	switch pKind {
	case cue.TopKind:
		return orig.Source(), nil

	// structs or objects, recurse and loop on fields
	case cue.StructKind:
		return pickStructs(indent + "  ", tag + "/structs", orig, pick)

	// lists or arrays, ... what do we do here?
	case cue.ListKind:
		return pickLists(indent + "  ", tag + "/lists", orig, pick)

	case cue.BoolKind, cue.NumberKind, cue.BytesKind, cue.StringKind:
		return pickBasicLit(indent + "  ", tag + "/basic", orig, pick)

	default:
		return val, fmt.Errorf("Unknown Cue Type %q in pickValues in type switch: %v %d", pKind, pick, pKind)
	}

	// fmt.Printf("%spick: %s  END\n", indent, tag)
	// return val, err
}

func pickStructs(indent, tag string, orig, pick cue.Value) (val ast.Node, err error) {

	//oKind := orig.Kind()
	//oLabel, _ := orig.Label()
	//pKind := pick.IncompleteKind()
	// pLabel, _ := pick.Label()

	// fmt.Printf("%sstructs: %q %q %q %q %q  BEG\n", indent, tag, oLabel, oKind, pLabel, pKind)

	fields := []ast.Decl{}
	oStruct, _ := orig.Struct()
	iter := oStruct.Fields()
	for iter.Next() {
		ol := iter.Label()
		ov := iter.Value()
		// ok := ov.Kind()

		// try to look it up
		pv := pick.Lookup(ol)
		// pk := pv.IncompleteKind()
		// TODO, did we find the kind?

		// TODO, do we want to condition based on pick being concrete or not?

		// print everything we decide on
		// fmt.Printf("%s  field: %s %s %s %s %s %d\n", indent, ol, ov, ok, pv, pk, len(fields))
		u := pv.Unify(ov)
		// fmt.Printf("%s  unify: %v %v %v %v %v\n%v %v %v %v\nunified: %v %v\n", indent,
			//ol, ok, ov.IsConcrete(), ov.IsClosed(), ov,
			//pk, pv.IsConcrete(), pv.IsClosed(), pv,
			//u, u.Err(),
		//)
		// fmt.Printf("%s    check: %v %v\n", indent, u.Equals(ov), u.Kind() == cue.BottomKind)

		// If there was an error during unification, continue to next field
		if u.Err() != nil {
			// fmt.Printf("%s    ERROR: %t %v\n", indent, err, err)
			continue
		}

		// if we unified, recurse
		rv, err := pickValues(indent + "  ", tag + "/field", ov, pv)
		if err != nil {
			// fmt.Printf("%s    ERROR: %t %v\n", indent, err, err)
			continue
		}

		fields = append(fields, rv.(ast.Decl))

		continue

		/*
		// If we unified
		if u.Equals(ov) {
			// fmt.Printf("%s    equals: %q %s\n", indent, ol, ov)
			syntax := ov.Syntax()
			field := &ast.Field {
				Label: ast.NewString(ol),
				Value: syntax.(ast.Expr),
				Token: token.COLON,
			}
			fields = append(fields, field)
			continue
		}

		if isBuiltin(pv) {
			// fmt.Println(indent, "   Builtin! ", ol, pk, pv)
		}

		// is it not found?
		if pv.Kind() == cue.BottomKind {
			// fmt.Printf("%s    continue on _|_ for %s\n", indent, ol)
			continue
		}

		switch pk {

		case pKind:
			// fmt.Printf("%s    samekind: %q %s\n", indent, ol, pv)
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
			// fmt.Printf("%s    default: %q %s\n", indent, ol, pv)
			syntax := pv.Syntax()
			field := &ast.Field {
				Label: ast.NewString(ol),
				Value: syntax.(ast.Expr),
				Token: token.COLON,
			}
			fields = append(fields, field)
		}
	*/
	}

	val = ast.NewStruct()
	s := val.(*ast.StructLit)
	fmt.Println()
	for _, f := range fields {
		s.Elts = append(s.Elts, f)
	}
	// fmt.Printf("%sstructs: %s  END %v %v\n", indent, tag, val, fields)
	return val, err
}

func pickLists(indent, tag string, orig, pick cue.Value) (val ast.Node, err error) {
	//oKind := orig.Kind()
	//oLabel, _ := orig.Label()
	//pKind := pick.Kind()
	//pLabel, _ := pick.Label()

	// fmt.Printf("%slists: %q %q %q %q %q  BEG\n", indent, tag, oLabel, oKind, pLabel, pKind)



	// fmt.Printf("%slists: %s  END\n", indent, tag)
	return val, err
}

// compares and picks basic lits
func pickBasicLit(indent, tag string, orig, pick cue.Value) (val ast.Node, err error) {
	// fmt.Println("BasicLit", orig, pick)
	u := pick.Unify(orig)

	if err := u.Err(); err != nil {
		return pick.Source(), err
	}

	return orig.Source(), err
}

func pickKinds(indent, tag string, orig, pick cue.Value) (val ast.Node, err error) {
	//oKind := orig.Kind()
	//oLabel, _ := orig.Label()
	pKind := pick.IncompleteKind()
	// pLabel, _ := pick.Label()

	// fmt.Printf("%skinds: %q %q %q %q %q  BEG\n", indent, tag, oLabel, oKind, pLabel, pKind)

	u := pick.Unify(orig)
	// fmt.Println("Unify", orig, pick, u)

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

	// fmt.Printf("%skinds: %s  END\n", indent, tag)
	return val, err
}
