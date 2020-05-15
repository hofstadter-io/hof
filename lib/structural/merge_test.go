package structural_test

import (
	"regexp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/hofstadter-io/hof/lib/structural"
)

type MergeTestSuite struct {
	suite.Suite

	mergeRT *structural.CueRuntime
}

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
