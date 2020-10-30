package structural_test

import (
	"fmt"
	"testing"

	"cuelang.org/go/cue"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/hofstadter-io/hof/lib/structural"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

var (
	PickFmtStr = "PickCases[%v]: %v"
	PickTestCases = []string{
		"#PickCases",
	}
)

type PickTestSuite struct {
	*cuetils.TestSuite
}

func NewPickTestSuite() *PickTestSuite {
	ts := cuetils.NewTestSuite(nil, PickOp)
	return &PickTestSuite{ ts }
}

func TestPickTestSuites(t *testing.T) {
	suite.Run(t, NewPickTestSuite())
}

func PickOp(name string, args cue.Value) (val cue.Value, err error) {
	orig := args.Lookup("orig")
	pick := args.Lookup("pick")
	return structural.PickValues(orig, pick)
}

func (PTS *PickTestSuite) TestPickCases() {

	err := PTS.SetupCue()
	assert.Nil(PTS.T(), err, fmt.Sprintf(PickFmtStr, "setup", "Loading test cases should return non-nil error"))
	if err != nil {
		return
	}

	PTS.Op = PickOp
	PTS.RunCases(PickTestCases)
}
