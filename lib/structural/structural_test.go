package structural_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/hofstadter-io/hof/lib/structural"
)

// This test function drives the top level suites for structural
func TestStructuralTestSuite(t *testing.T) {

	suite.Run(t, new(CommonTestSuite))
	suite.Run(t, new(DiffTestSuite))
	suite.Run(t, new(MergeTestSuite))
	suite.Run(t, new(MaskTestSuite))
	suite.Run(t, new(PickTestSuite))
	suite.Run(t, new(QueryTestSuite))
	suite.Run(t, new(FindAttrsTestSuite))

}

func loadCueTestData(entrypoints ...string) *structural.CueRuntime {
	cr := structural.NewCueRuntime()
	errs := cr.LoadCue(entrypoints)
	if len(errs) > 0 {
		for _, e := range errs {
			fmt.Println(e)
		}
		panic("Errors loading cue test data for entrypoints: " + strings.Join(entrypoints, ", "))
	}

	return cr
}
