package structural_test

import (
	"regexp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/hofstadter-io/hof/lib/structural"
)

type MaskTestSuite struct {
	suite.Suite

	maskRT *structural.CueRuntime
}

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
