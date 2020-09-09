package test

import (
	"time"

	"cuelang.org/go/cue"

	"github.com/hofstadter-io/hof/lib/cuetils"
)

// A suite is a collection of testers
type Suite struct {
	// Name of the Suite (Cue field)
	Name string

	// Cue Value for the Suite
	Value cue.Value

	// Extracted testers, including any selectors
	Tests  []Tester

	// pass/fail/skip stats
	Stats Stats

	// Total Suite runtime (to account for gaps between tests)
	Runtime time.Duration

	// Errors encountered during the testing
	Errors []error
}

func getValueTestSuites(val cue.Value, labels []string) ([]Suite, error) {
	vals, err := cuetils.GetByAttrKeys(val, "test", append(labels, "suite"), nil)
	suites := []Suite{}
	for _, v := range vals {
		suites = append(suites, Suite{Name: v.Key, Value: v.Val})
	}
	return suites, err
}

// A tester has configuration for running a set of tests
type Tester struct {
	// Name of the Tester (Cue field)
	Name string

	// Type of the Tester (@test(key[0]))
	Type string

	// Cue Value for the Tester
	Value cue.Value

	// Execution output
	Output string

	// pass/fail/skip stats
	Stats Stats

	// Errors encountered during the testing
	Errors []error
}

func getValueTestSuiteTesters(val cue.Value, labels []string) ([]Tester, error) {
	vals, err := cuetils.GetByAttrKeys(val, "test", labels, []string{})
	testers := []Tester{}
	for _, v := range vals {
		a := v.Val.Attribute("test")
		typ, err := a.String(0)
		if err != nil {
			return testers, err
		}
		testers = append(testers, Tester{Name: v.Key, Type: typ, Value: v.Val})
	}
	return testers, err
}

type Stats struct {
	Pass int
	Fail int
	Skip int

	Start time.Time
	End   time.Time
	Time  time.Duration
}

func (S *Stats) add(s Stats) {
	S.Pass += s.Pass
	S.Fail += s.Fail
	S.Skip += s.Skip
	S.Time += s.Time
}
