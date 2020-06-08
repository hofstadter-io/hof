package test

import (
	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

func RunTestFromArgsFlags(args []string, cmdflags flags.TestFlagpole) (error) {

	// Loadup our Cue files
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

	// Is the user only looking for information
	if cmdflags.List {
		printTests(suites, false)
		return nil
	}

	// Run all of our suites
	_, err = RunSuites(suites, -1)
	if err != nil {
		return err
	}

	// TODO, print errors

	// Print our final tests and stats
	printTests(suites, true)

	return nil
}

