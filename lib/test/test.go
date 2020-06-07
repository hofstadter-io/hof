package test

import (
	"fmt"
	"regexp"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

type Test struct {
	Name string
	Type string
	Value cue.Value

	// for nesting
	Tests []Test
}


func RunTestFromArgsFlags(args []string, cmdflags flags.TestFlagpole) (error) {

	// fmt.Printf("Test: %v %#+v\n", args, cmdflags)

	crt, err := cuetils.CueRuntimeFromEntrypointsAndFlags(args)
	if err != nil {
		return err
	}

	err = crt.Load()
	if err != nil {
		return err
	}

	// Get test suites from top level
	suites, err := getValueTestSuites(crt.CueValue, "ALL", cmdflags.Suite)
	if err != nil {
		return err
	}

	// find tests in suites
	for s, suite := range suites {
		ts, err := getValueTestSuiteTesters(suite.Value, suite.Name, cmdflags.Tester)
		if err != nil {
			return err
		}
		// make sure to write to original
		suites[s].Tests = ts
	}

	if cmdflags.List {
		printTests("", suites)
	}

	return nil

}

func printTests(prefix string, tests []Test) {
	if tests == nil || len(tests) == 0 {
		return
	}
	for _, T := range tests {
		as := []string{}
		for _, A := range T.Value.Attributes() {
			vals := A.Vals()
			a := fmt.Sprintf("%s:%v", A.Name(), vals)
			as = append(as, a)
		}

		fmt.Printf(
			"%-16s  %-16s  [%d] %v\n",
			prefix + "[" + T.Type + "]",
			T.Name,
			len(T.Tests),
			as,
		)

		// recurse
		printTests(prefix + "  ", T.Tests)
		if prefix == "" {
			fmt.Println()
		}
	}

}

func getValueTestSuites(val cue.Value, name string, labels []string) ([]Test, error) {
	tests, err := getValueByAttrKeys(val, name, "test", []string{"suite"}, labels)
	for i, _ := range tests {
		// make sure to write to original
		tests[i].Type = "suite"
	}
	return tests, err
}

func getValueTestSuiteTesters(val cue.Value, name string, labels []string) ([]Test, error) {
	tests, err := getValueByAttrKeys(val, name, "test", []string{}, labels)
	for i, _ := range tests {
		// make sure to write to original
		tests[i].Type = "tester"
	}
	return tests, err
}

// Todo, rewrite this to use structural
func getValueByAttrKeys(val cue.Value, name, attr string, all, any []string) ([]Test, error) {
	// fmt.Println("GET:", name, attr, all, any)
	rets := []Test{}

	S, err := val.Struct()
	if err != nil {
		es := errors.Errors(err)
		for _, e := range es {
			fmt.Println(e)
		}
		return rets, fmt.Errorf("Error loading cue code")
	}

	// Loop through all top level fields
	iter := S.Fields()
	for iter.Next() {

		label := iter.Label()
		value := iter.Value()
		attrs := value.Attributes()

		// fmt.Println("  -", label, attrs)

		// find top-level with gen attr
		hasattr := false
		for _, A := range attrs {
			// does it have an "@<attr>(...)"
			if A.Name() ==  attr {
				vals := A.Vals()

				// must match all
				if len(all) > 0 {
					match := true

					// loop over the all list
					for _, l := range all {
						R := regexp.MustCompile(l)
						// loop over the field attt key names
						found := false
						for v, _  := range vals {
							m := R.MatchString(v)
							if m {
								found = true
								break
							}
						}
						// break one more time if we have failed
						if !found {
							match = false
							break
						}
					}

					// did we not match all?
					if !match {
						continue
					}
				}

				// match one of any
				if len(any) > 0 {
					match := false

					// loop over the any list
					for _, l := range any {
						R := regexp.MustCompile(l)
						// loop over the field attt key names
						for v, _  := range vals {
							m := R.MatchString(v)
							if m {
								match = true
								break
							}
						}

						// break again if we have matched
						if match {
							break
						}
					}

					// did we not match any?
					if !match {
						continue
					}
				}

				// fmt.Println("  ...Has", label, A.Name())
				// passed, we should include
				hasattr = true
				break
			}
		}

		// fmt.Println("  ...Attr", label, attr, hasattr)
		// ok, we're back outside the attrs look now, did we match on it?
		// if no, let's try the next field
		if !hasattr {
			continue
		}

		// add it and move on!
		rets = append(rets, Test{Name: label, Value: value})
	}

	return rets, nil
}
