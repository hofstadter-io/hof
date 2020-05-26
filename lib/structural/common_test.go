package structural_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/hofstadter-io/hof/lib/structural"
)

type CommonTestSuite struct {
	suite.Suite
}

func TestCommonTestSuites(t *testing.T) {
	suite.Run(t, new(CommonTestSuite))
}

func (suit *CommonTestSuite) TestPV() {
	pv := structural.NewpvStruct()
	pv.Ensure("ok")
	ok := pv.Get("ok")
	assert.NotNil(suit.T(), ok)

	ok.Set("blah", *structural.NewpvStruct().ToExpr())

	ok = ok.Get("blah")
	assert.NotNil(suit.T(), ok)

	assert.NotNil(suit.T(), pv.Get("ok").Get("blah"))
}
