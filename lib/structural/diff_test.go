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
	DiffFmtStr = "DiffCases[%v]: %v"
	DiffTestCases = []string{
		"#DiffCases",
	}
)

type DiffTestSuite struct {
	*cuetils.TestSuite
}

func NewDiffTestSuite() *DiffTestSuite {
	ts := cuetils.NewTestSuite(nil, DiffOp)
	return &DiffTestSuite{ ts }
}

func TestDiffTestSuites(t *testing.T) {
	suite.Run(t, NewDiffTestSuite())
}

func DiffOp(name string, args cue.Value) (val cue.Value, err error) {
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

	next := args.Lookup("next")
	pSyn, pErr := cuetils.ValueToSyntaxString(next)
	if pErr != nil {
		fmt.Println(pSyn)
		return val, pErr
	}

	return structural.DiffValues(orig, next)
}

func (PTS *DiffTestSuite) TestDiffCases() {

	err := PTS.SetupCue()
	assert.Nil(PTS.T(), err, fmt.Sprintf(DiffFmtStr, "setup", "Loading test cases should return non-nil error"))
	if err != nil {
		return
	}

	tSyn, err := cuetils.ValueToSyntaxString(PTS.CRT.CueValue)
	assert.Nil(PTS.T(), err, fmt.Sprintf(DiffFmtStr, "syntax", "Printing test cases should return non-nil error"))
	if err != nil {
		fmt.Println(tSyn)
		return
	}

	PTS.Op = DiffOp
	PTS.RunCases(DiffTestCases)
}

/*
func (suit *DiffTestSuite) TestDiff() {
	result, err := structural.CueDiff("{a:1, c:[1,2], x:1, s: {ssss: 2, ss: 2}}", "{b:2, c:[1,2,3], x:2, s: {ss: 1, sss: 1}}")
	assert.Nil(suit.T(), err)
	expected := `changed: {
        x: {
                from: 1
                to:   2
        }
        c: {
                from: [1, 2]
                to: [1, 2, 3]
        }
}
removed: {
        a: 1
}
inplace: {
        s: {
                changed: {
                        ss: {
                                from: 2
                                to:   1
                        }
                }
                removed: {
                        ssss: 2
                }
                added: {
                        sss: 1
                }
        }
}
added: {
        b: 2
}`
	space := regexp.MustCompile(`\s+`)
	result = space.ReplaceAllString(result, " ")
	expected = space.ReplaceAllString(expected, " ")
	assert.Equal(suit.T(), result, expected)
}
*/
