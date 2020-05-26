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
	MaskFmtStr = "MaskCases[%v]: %v"
	MaskTestCases = []string{
		"#MaskCases",
	}
)

type MaskTestSuite struct {
	*cuetils.TestSuite
}

func NewMaskTestSuite() *MaskTestSuite {
	ts := cuetils.NewTestSuite(nil, MaskOp)
	return &MaskTestSuite{ ts }
}

func TestMaskTestSuites(t *testing.T) {
	suite.Run(t, NewMaskTestSuite())
}

func MaskOp(name string, args cue.Value) (val cue.Value, err error) {
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

	mask := args.Lookup("mask")
	pSyn, pErr := cuetils.ValueToSyntaxString(mask)
	if pErr != nil {
		fmt.Println(pSyn)
		return val, pErr
	}

	return structural.MaskValues(orig, mask)
}

func (PTS *MaskTestSuite) TestMaskCases() {

	err := PTS.SetupCue()
	assert.Nil(PTS.T(), err, fmt.Sprintf(MaskFmtStr, "setup", "Loading test cases should return non-nil error"))
	if err != nil {
		return
	}

	tSyn, err := cuetils.ValueToSyntaxString(PTS.CRT.CueValue)
	assert.Nil(PTS.T(), err, fmt.Sprintf(MaskFmtStr, "syntax", "Printing test cases should return non-nil error"))
	if err != nil {
		fmt.Println(tSyn)
		return
	}

	PTS.Op = MaskOp
	PTS.RunCases(MaskTestCases)
}

/*
func (suit *MaskTestSuite) TestMask() {
	result, err := structural.CueMask("{a:1, cc: [1,2,3], c:[1,2,3], x:1, s: {ssss: 2, ss: 2}}", "{a: int, cc: [1,1,1], c: <3, s: {a: string, ss: >2}}")
	assert.Nil(suit.T(), err)
	expected := `x: 1
c: [3]
s: {
        ssss: 2
        ss:   2
}
cc: [2, 3]`

	space := regexp.MustCompile(`\s+`)
	result = space.ReplaceAllString(result, " ")
	expected = space.ReplaceAllString(expected, " ")
	assert.Equal(suit.T(), result, expected)
}
*/
