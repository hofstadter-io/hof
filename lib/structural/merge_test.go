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
	MergeFmtStr = "MergeCases[%v]: %v"
	MergeTestCases = []string{
		"#MergeCases",
	}
)

type MergeTestSuite struct {
	*cuetils.TestSuite
}

func NewMergeTestSuite() *MergeTestSuite {
	ts := cuetils.NewTestSuite(nil, MergeOp)
	return &MergeTestSuite{ ts }
}

func TestMergeTestSuites(t *testing.T) {
	suite.Run(t, NewMergeTestSuite())
}

func MergeOp(name string, args cue.Value) (val cue.Value, err error) {
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

	merge := args.Lookup("merge")
	pSyn, pErr := cuetils.ValueToSyntaxString(merge)
	if pErr != nil {
		fmt.Println(pSyn)
		return val, pErr
	}

	return structural.MergeValues(orig, merge)
}

func (PTS *MergeTestSuite) TestMergeCases() {

	err := PTS.SetupCue()
	assert.Nil(PTS.T(), err, fmt.Sprintf(MergeFmtStr, "setup", "Loading test cases should return non-nil error"))
	if err != nil {
		return
	}

	tSyn, err := cuetils.ValueToSyntaxString(PTS.CRT.CueValue)
	assert.Nil(PTS.T(), err, fmt.Sprintf(MergeFmtStr, "syntax", "Printing test cases should return non-nil error"))
	if err != nil {
		fmt.Println(tSyn)
		return
	}

	PTS.Op = MergeOp
	PTS.RunCases(MergeTestCases)
}
/*
func (suit *MergeTestSuite) TestMerge() {
	result, err := structural.CueMerge("{a: 1, c: 2, d: {c:1,b:2}}", "{b : 2, c:3, d: {a: 1,d:2}}")
	assert.Nil(suit.T(), err)
	expected := `a: 1
c: 3
b: 2
d: {
        a: 1
        c: 1
        b: 2
        d: 2
}`

	space := regexp.MustCompile(`\s+`)
	result = space.ReplaceAllString(result, " ")
	expected = space.ReplaceAllString(expected, " ")
	assert.Equal(suit.T(), result, expected)
}
*/
