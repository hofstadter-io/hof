package cmds

import (
	"github.com/hofstadter-io/hofmod-cli/schema"
)

#TestCommand: schema.#Command & {
	TBD:   "α"
	Name:  "test"
	Usage: "test"
	Aliases: ["t"]
	Short: "test all sorts of things"
	Long:  #TestCommandHelp
}

#TestCommandHelp: #"""
api:    test rest and graphqlendpoints
bdd:    behaviour style tests
table:  table based tests
script: testscript based tests
story:  user story based tests
bench:  benchmark based tests
chaos:  chaos testing
e2e:    end-to-end integration tests
suite:  nested and grouped sets of tests and other suites
"""#
