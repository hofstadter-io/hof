package structural_test

import (
	"regexp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/hofstadter-io/hof/lib/structural"
)

type QueryTestSuite struct {
	suite.Suite

	queryRT *structural.CueRuntime
}

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
