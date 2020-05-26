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
	aSyn, aErr := cuetils.ValueToSyntaxString(args)
	if aErr != nil {
		fmt.Println(aSyn)
		return val, aErr
	}

	orig := args.Lookup("orig")
	oSyn, oErr := cuetils.ValueToSyntaxString(orig)
	if oErr != nil {
		fmt.Println(oSyn)
		return val, oErr
	}

	pick := args.Lookup("pick")
	pSyn, pErr := cuetils.ValueToSyntaxString(pick)
	if pErr != nil {
		fmt.Println(pSyn)
		return val, pErr
	}

	return structural.PickValues(orig, pick)
}

func (PTS *PickTestSuite) TestPickCases() {

	err := PTS.SetupCue()
	assert.Nil(PTS.T(), err, fmt.Sprintf(PickFmtStr, "setup", "Loading test cases should return non-nil error"))
	if err != nil {
		return
	}

	tSyn, err := cuetils.ValueToSyntaxString(PTS.CRT.CueValue)
	assert.Nil(PTS.T(), err, fmt.Sprintf(PickFmtStr, "syntax", "Printing test cases should return non-nil error"))
	if err != nil {
		fmt.Println(tSyn)
		return
	}

	PTS.Op = PickOp
	PTS.RunCases(PickTestCases)
}
