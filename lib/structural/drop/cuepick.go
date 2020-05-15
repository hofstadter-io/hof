package main

import (
	"fmt"
	"os"
	"reflect"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/ast/astutil"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/load"
)

var (
	Entrypoints = []string{"./"}

	FieldOpts = []cue.Option{
		cue.Attributes(true),
		cue.Concrete(false),
		cue.Definitions(true),
		cue.Docs(true),
		cue.Hidden(true),
		cue.Optional(true),
	}

	CueRT *cue.Runtime
	BIS   []*build.Instance

	CueInstance    *cue.Instance
	TopLevelValue  cue.Value
	TopLevelStruct *cue.Struct

	TransformedValue  cue.Value
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {

	fmt.Println("Starting....")
	fmt.Println("==========")

	errs := loadCue(Entrypoints)
	if len(errs) > 0 {
		fmt.Println("Errors:", errs)
		return
	}

	var err error

	//fmt.Println("========== top-level structs")
	//err = printCueStructs()
	//checkErr(err)

	fmt.Println("========== top-level values")
	err = printCueValues()
	checkErr(err)

  fmt.Println("========== xxx-transforming-xxx")
	err = transformValue(TopLevelValue.Lookup("A"), "")
	checkErr(err)

	fmt.Println("========== transformation")
	err = printValue(TopLevelValue.Lookup("A"))

  //fmt.Println("---------- xxx-picking-xxx")
	//TransformedValue, err = Pick(TopLevelValue.Lookup("A"), TopLevelValue.Lookup("P"))
	//checkErr(err)

  fmt.Println("---------- new value")
	err = printValue(TransformedValue)
	checkErr(err)

	fmt.Println("========== done")
}


func loadCue(entrypoints []string) []error {
	var errs []error

	CueRT = &cue.Runtime{}
	BIS = load.Instances(entrypoints, &load.Config{
		Package: "",
		Tools: false,
	})

	fmt.Println("len(BIS):", len(BIS))

	// Build the Instances
	I, err := CueRT.Build(BIS[0])
	if err != nil {
		es := errors.Errors(err)
		// fmt.Println("BUILD ERR", es, I)
		for _, e := range es {
			errs = append(errs, e.(error))
		}
		return errs
	}

	CueInstance = I

	// Get top level value from cuelang
	V := I.Value()
	TopLevelValue = V

	// Get top level struct from cuelang
	S, err := V.Struct()
	if err != nil {
		// fmt.Println("STRUCT ERR", err)
		es := errors.Errors(err)
		for _, e := range es {
			errs = append(errs, e.(error))
		}
		return errs
	}

	TopLevelStruct = S

	return errs
}

func printCueValues() error {
	// Loop through all top level fields
	iter, err := TopLevelValue.Fields(FieldOpts...)
	if err != nil {
		return err
	}

	for iter.Next() {

		value := iter.Value()

		walkValue(value, "")
	}

	return nil
}

func walkValue(val cue.Value, indent string) error {
	L, _ := val.Label()
	fmt.Println(indent, L, val.Attributes())

	iter, err := val.Fields(FieldOpts...)
	if err != nil {
		return err
	}

	for iter.Next() {

		label := iter.Label()
		value := iter.Value()

		if value.Kind() == cue.StructKind {
			walkValue(value, indent + "  ")
		} else {
			fmt.Println(indent, "--", label, value.Kind(), value.Attributes())
		}

	}

	fmt.Println()

	return nil
}

func printCueStructs() error {
	// Loop through all top level fields
	iter := TopLevelStruct.Fields(FieldOpts...)

	for iter.Next() {

		label := iter.Label()
		value := iter.Value()

		fmt.Println(label, value, value.Kind(), value.Attributes())
	}

	return nil
}

func printValue(val cue.Value) error {
	node := val.Syntax(FieldOpts...)

	bytes, err := format.Node(
		node,
		format.TabIndent(false),
		format.UseSpaces(2),
	)
	if err != nil {
		return err
	}

	fmt.Println(string(bytes))

	return nil
}

func transformValue(val cue.Value, indent string) error {
	L, _ := val.Label()
	fmt.Println(indent, L, val.Attributes())

	iter, err := val.Fields(FieldOpts...)
	if err != nil {
		return err
	}

	for iter.Next() {

		label := iter.Label()
		value := iter.Value()

		if value.Kind() == cue.StructKind {
			transformValue(value, indent + "  ")
		} else {
			fmt.Println(indent, "--", label, value.Kind(), value.Attributes())
		}

	}

	fmt.Println()

	if indent == "" {
		i1 := ast.NewIdent("newval")
		e1 := ast.NewString("hello")
		f1 := &ast.Field{
			Label: i1,
			Value: e1,
		}

		i2 := ast.NewIdent("moo")
		e2 := ast.NewString("bar")
		as2 := &ast.Attribute{Text: "@expanded(BIG)"}
		f2 := &ast.Field{
			Label: i2,
			Value: e2,
			Attrs: []*ast.Attribute{as2},
		}


		s := ast.NewStruct(f1, f2)

		node := val.Syntax(FieldOpts...)

		before := func(c astutil.Cursor) bool {
			switch N := c.Node().(type) {
				case *ast.Field:
					if len(N.Attrs) > 0 {
						if N.Attrs[0].Text == "@delete()" {
							c.Delete()
						}
						if N.Attrs[0].Text == "@expand()" {
							N.Value = s
							N.Attrs = nil
							// c.Replace(f1)
						}
					}
			}
			return true
		}

		n := astutil.Apply(node, before, nil)

		// b := TopLevelValue.Lookup("B")
		// b.Fill(s)
		v := CueInstance.Eval(n.(*ast.StructLit))
		TransformedValue = v
	}

	return nil
}

func Pick(orig, pick cue.Value) (cue.Value, error) {

	o := orig.Syntax(FieldOpts...)
	p := pick.Syntax(FieldOpts...)

	n, err := pickNode(o, p)
	if err != nil {
		return cue.Value{}, err
	}

	v := CueInstance.Eval(n.(*ast.StructLit))
	return v, nil
}

func pickNode(orig, pick ast.Node) (ast.Node, error) {

	var errs = []error{}

	before := func(c astutil.Cursor) bool {
		switch N := c.Node().(type) {
			case *ast.Field:
				fmt.Printf("Field -- %v: %v\n", N.Label, N.Attrs)
				// Just check that they unify instead?
				O, ook := orig.(*ast.StructLit)
				_, pok := pick.(*ast.StructLit)

				if ook && pok {

					fmt.Println("len(O.Elts) =", len(O.Elts))

					for _, el := range O.Elts {

						fmt.Println(" -", el, reflect.TypeOf(el))
						f, ok := el.(*ast.Field)
						if !ok {
							continue
						}

						fmt.Println("   + FIELD", f.Label)
						nl, ni, ne := ast.LabelName(N.Label)
						fl, fi, fe := ast.LabelName(f.Label)

						if fl == nl {
							fmt.Println("   + MATCH", nl, fl, ni, fi, ne, fe)
							c.Replace(f)
							break
						}
					}

				}

			default:
				err := fmt.Errorf("UNHANDLED node type in before() '%v'", reflect.TypeOf(c.Node()))
				errs = append(errs, err)
				// return false
		}
		return true
	}

	n := astutil.Apply(pick, before, nil)

	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
	}

	return n, nil
}

