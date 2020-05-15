package structural_test

import (
	"regexp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/hofstadter-io/hof/lib/structural"
)

type PickTestSuite struct {
	suite.Suite

	pickRT *structural.CueRuntime
}

func (suit *PickTestSuite) TestPick() {
	result, err := structural.CuePick("{a:1, cc: [1,2,3], c:[1,2], x:1, s: {ssss: 2, ss: 2}}", "{a: int, cc: [1,1,1], c: <3, s: {a: int, ss: <5}}")
	assert.Nil(suit.T(), err)
	expected := `a: 1
c: [1, 2]
s: {
        ss: 2
}
cc: [1]`

	space := regexp.MustCompile(`\s+`)
	result = space.ReplaceAllString(result, " ")
	expected = space.ReplaceAllString(expected, " ")
	assert.Equal(suit.T(), result, expected)
}
