package structural_test

import (
	"fmt"
	"testing"

	"cuelang.org/go/cue"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/suite"

	"github.com/hofstadter-io/hof/lib/structural"
	"github.com/hofstadter-io/hof/lib/cuetils"
)

var (
	QueryFmtStr = "QueryCases[%v]: %v"
	QueryTestCases = []string{
		"#QueryCases",
	}
)

type QueryTestSuite struct {
	*cuetils.TestSuite
}

func NewQueryTestSuite() *QueryTestSuite {
	ts := cuetils.NewTestSuite(nil, QueryOp)
	return &QueryTestSuite{ ts }
}

func TestQueryTestSuites(t *testing.T) {
	// suite.Run(t, NewQueryTestSuite())
}

func QueryOp(name string, args cue.Value) (val cue.Value, err error) {
	aSyn, aErr := cuetils.PrintCueValue(args)
	if aErr != nil {
		fmt.Println(aSyn)
		return val, aErr
	}

	orig := args.Lookup("orig")
	oSyn, oErr := cuetils.PrintCueValue(orig)
	if oErr != nil {
		fmt.Println(oSyn)
		return val, oErr
	}

	query := args.Lookup("query")
	pSyn, pErr := cuetils.PrintCueValue(query)
	if pErr != nil {
		fmt.Println(pSyn)
		return val, pErr
	}

	return structural.QueryValues(orig, query)
}

func (PTS *QueryTestSuite) TestQueryCases() {

	err := PTS.SetupCue()
	assert.Nil(PTS.T(), err, fmt.Sprintf(QueryFmtStr, "setup", "Loading test cases should return non-nil error"))
	if err != nil {
		return
	}

	tSyn, err := cuetils.ValueToSyntaxString(PTS.CRT.CueValue)
	assert.Nil(PTS.T(), err, fmt.Sprintf(QueryFmtStr, "syntax", "Printing test cases should return non-nil error"))
	if err != nil {
		fmt.Println(tSyn)
		return
	}

	PTS.Op = QueryOp
	PTS.RunCases(QueryTestCases)
}
/*
func (suit *QueryTestSuite) TestQuery() {
	input := "{b : 2, c:3, d: {a: 1, d:2, x: \"hi\"}}"
	result, err := structural.CueQuery("[int]", input)
	assert.Nil(suit.T(), err)
	expected := `[3, 2]`

	space := regexp.MustCompile(`\s+`)
	result = space.ReplaceAllString(result, " ")
	expected = space.ReplaceAllString(expected, " ")
	assert.Equal(suit.T(), result, expected)

	result, err = structural.CueQuery("[_]", input)
	assert.Nil(suit.T(), err)
	expected = `[3, 2, {
x: "hi"
a: 1
d: 2
}]`

	space = regexp.MustCompile(`\s+`)
	result = space.ReplaceAllString(result, " ")
	expected = space.ReplaceAllString(expected, " ")
	assert.Equal(suit.T(), result, expected)

	result, err = structural.CueQuery("[_,string]", input)
	assert.Nil(suit.T(), err)
	expected = `["hi"]`

	space = regexp.MustCompile(`\s+`)
	result = space.ReplaceAllString(result, " ")
	expected = space.ReplaceAllString(expected, " ")
	assert.Equal(suit.T(), result, expected)
}
*/
