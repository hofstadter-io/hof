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
	AttrsFmtStr = "AttrsCases[%v]: %v"
	AttrsTestCases = []string{
		"#AttrsCases",
	}
)

type AttrsTestSuite struct {
	*cuetils.TestSuite
}

func NewAttrsTestSuite() *AttrsTestSuite {
	ts := cuetils.NewTestSuite(nil, AttrsOp)
	return &AttrsTestSuite{ ts }
}

func TestAttrsTestSuites(t *testing.T) {
	// suite.Run(t, NewAttrsTestSuite())
}

func AttrsOp(name string, args cue.Value) (val cue.Value, err error) {
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

	aVal := args.Lookup("attrs")
	aValK := args.Lookup("attrsK")
	aValKV := args.Lookup("attrsKV")

	var (
		attrs   []string
		attrsK  map[string][]string
		attrsKV map[string]map[string]string
	)

	if aVal.Kind() != cue.BottomKind {
		err = aVal.Decode(attrs)
		if err != nil {
			return val, err
		}
	}
	if aValK.Kind() != cue.BottomKind {
		err = aValK.Decode(attrsK)
		if err != nil {
			return val, err
		}
	}
	if aValKV.Kind() != cue.BottomKind {
		err = aValKV.Decode(attrsKV)
		if err != nil {
			return val, err
		}
	}

	vals, err := structural.FindByAttrs(orig, attrs, attrsK, attrsKV)
	if err != nil {
		return val, err
	}

	return vals[0], nil
}

func (PTS *AttrsTestSuite) TestAttrsCases() {

	err := PTS.SetupCue()
	assert.Nil(PTS.T(), err, fmt.Sprintf(AttrsFmtStr, "setup", "Loading test cases should return non-nil error"))
	if err != nil {
		return
	}

	tSyn, err := cuetils.ValueToSyntaxString(PTS.CRT.CueValue)
	assert.Nil(PTS.T(), err, fmt.Sprintf(AttrsFmtStr, "syntax", "Printing test cases should return non-nil error"))
	if err != nil {
		fmt.Println(tSyn)
		return
	}

	PTS.Op = AttrsOp
	PTS.RunCases(AttrsTestCases)
}

/*
func (suit *FindAttrsTestSuite) TestFindAttrs() {
	var r cue.Runtime
	input := `
a: 2
b: _ @findme(hide=true)        // 1
b: 2
c: 2
d: _ @findme(hide=false)       // 2
d: 3
x: _ @andme()                  // 3
xx: _ @check(this=isthis)      // 4
xx1: _ @check(a,this=isthis)   // 5
xxy: _ @check(this=notthat)
xxz: _ @check(this)
abc: _ @hasthis(needs,these)   // 6
abd: _ @hasthis(a,needs,these) // 7
abb: _ @hasthis(eeds,hese)
`
	i, err := r.Compile("", input)
	assert.Nil(suit.T(), err)
	v := i.Value()
	assert.Nil(suit.T(), v.Err())
	vs, err := structural.FindByAttrs(v,
		[]string{"findme", "andme"},
		map[string][]string{
			"hasthis": []string{"needs", "these"},
		},
		map[string]map[string]string{
			"check": map[string]string{
				"this": "isthis",
			},
		})
	assert.Nil(suit.T(), err)
	assert.Len(suit.T(), vs, 7)
}
*/
