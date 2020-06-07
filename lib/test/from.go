package test

import (
	"fmt"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

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
	suites, err := getValueTestSuites(crt.CueValue, "suites", cmdflags.Suite)
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
		printTests(suites, false)
		return nil
	}

	TS, err := RunSuites(suites, -1)
	if err != nil {
		return err
	}


	printTests(suites, true)
	fmt.Printf("Final TS: %#+v\n", TS)
	return nil

}

