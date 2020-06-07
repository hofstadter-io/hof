package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/lib/test"
)

var testLong = `hof test - \(#TestCommand.Short)

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


`

func init() {

	TestCmd.Flags().BoolVarP(&(flags.TestFlags.List), "list", "", false, "list matching tests that would run")
	TestCmd.Flags().StringSliceVarP(&(flags.TestFlags.Suite), "suite", "s", nil, "<name>: _ @test(suite)'s to run")
	TestCmd.Flags().StringSliceVarP(&(flags.TestFlags.Tester), "tester", "t", nil, "<name>: _ @test(<tester>)'s to run")
	TestCmd.Flags().StringSliceVarP(&(flags.TestFlags.Environment), "env", "e", nil, "environment")
}

func TestRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = test.RunTestFromArgsFlags(args, flags.TestFlags)

	return err
}

var TestCmd = &cobra.Command{

	Use: "test",

	Aliases: []string{
		"t",
	},

	Short: "test all sorts of things",

	Long: testLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = TestRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := TestCmd.HelpFunc()
	usage := TestCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	TestCmd.SetHelpFunc(thelp)
	TestCmd.SetUsageFunc(tusage)

}