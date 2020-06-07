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

	Flags: [
		{
			Name:    "list"
			Type:    "bool"
			Default: "false"
			Help:    "list matching tests that would run"
			Long:    "list"
			Short:   ""
			...
		},
		{
			Name:    "suite"
			Type:    "[]string"
			Default: "nil"
			Help:    "<name>: _ @test(suite)'s to run"
			Long:    "suite"
			Short:   "s"
			...
		},
		{
			Name:    "tester"
			Type:    "[]string"
			Default: "nil"
			Help:    "<name>: _ @test(<tester>)'s to run"
			Long:    "tester"
			Short:   "t"
			...
		},
		{
			Name:    "environment"
			Type:    "[]string"
			Default: "nil"
			Help:    "environment"
			Long:    "env"
			Short:   "e"
			...
		},
	]

}

#TestCommandHelp: #"""
hof test - \(#TestCommand.Short)

hof test helps you test all the things by providing
a top-level driver and sitting on top of any tool.
You can group tests into Suites, nest and label them
and later run only the tests you want. Several builtin
Testers are available and patterns for testing your
applications, top to bottom and end to end.

Suites are a top level grouping attribute. You may go 
two levels deep for now, however there are both <name>
globs and labels to match with.

Here is an example "test.cue" file (testers are omitted):

------------------------------------------
MySuite: _ @test(suite)
MySuite: {

	// These sets will have nested testers, more on that below

	Unit: _ @test(suite,labelA,labelB)
	Unit: { ... }

	Regressions: _ @test(suite,labelA,frontend,backend)
	Regressions: { ... }

	"integration/frontend": _ @test(suite,frontend)
	"integration/frontend": {...}
	"integration/backend": _ @test(set,backend)
	"integration/backend": {...}
}

// These could have nested suites themselves
"service-f/fast": _ @test(suite,frontend)
"service-f/fast": {...}
"service-f/slow": _ @test(suite,frontend)
"service-f/slow": {...}

"service-b/fast": _ @test(suite,backend)
"service-b/fast": {...}
"service-b/slow": _ @test(suite,backend)
"service-b/slow": {...}
------------------------------------------


Testers are the pieces that run actual tests. They may:

- delegate by running another program or script
- use many of the builtin, generic testers
- use one of the special purpose systems built into hof  

All tester implementations can be found under the "lib/test" directory.

Testers:

	script:        step based testing that can work with the terminal, http, and data
	table:         table based testing using cue to express more cases with fewer lines
	http:          test rest and graphql endpoints, use "hof test import" to get a jump start
	exec:          exec out to the shell to run bash or anything available in your environment
	story:         behaviour style tests, in the syntax of Gherkin
	bench:         benchmark based tests, often built from other testers



"""#
