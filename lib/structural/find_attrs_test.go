package structural_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/hofstadter-io/hof/lib/structural"

	"cuelang.org/go/cue"
)

type FindAttrsTestSuite struct {
	suite.Suite

	findAttrsRT *structural.CueRuntime
}

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
