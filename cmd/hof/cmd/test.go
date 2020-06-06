package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var testLong = `api:    test rest and graphqlendpoints
bdd:    behaviour style tests
table:  table based tests
script: testscript based tests
story:  user story based tests
bench:  benchmark based tests
chaos:  chaos testing
e2e:    end-to-end integration tests
suite:  nested and grouped sets of tests and other suites`

func TestRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

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
