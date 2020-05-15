package structural_test

import (
	"regexp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/hofstadter-io/hof/lib/structural"
)

type DiffTestSuite struct {
	suite.Suite

	diffRT *structural.CueRuntime
}

func (suite *DiffTestSuite) SetupTest() {
	suite.diffRT = loadCueTestData("testdata/diff.cue")
}

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
