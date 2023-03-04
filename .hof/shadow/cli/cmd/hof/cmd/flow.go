package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var flowLong = `run CUE pipelines with the hof/flow DAG engine

Use hof/flow to transform data, call APIs, work with DBs,
read and write files, call any program, handle events,
and much more.

'hof flow' is very similar to 'cue cmd' and built on the same flow engine.
Tasks and dependencies are inferred.
Hof flow has a slightly different interface and more task types.

Docs: https://docs.hofstadter.io/data-flow

Example:

  @flow()

  call: {
    @task(api.Call)
    req: { ... }
    resp: {
      statusCode: 200
      body: string
    }
  }

  print: {
    @task(os.Stdout)
    test: call.resp
  }

Arguments:
  cue entrypoints are the same as the cue cli
  @path/name  is shorthand for -f / --flow should match the @flow(path/name)
  +key=value  is shorthand for -t / --tags and are the same as CUE injection tags

  arguments can be in any order and mixed

@flow() indicates a flow entrypoint
  you can have many in a file or nested values
  you can run one or many with the -f flag

@task() represents a unit of work in the flow dag
  intertask dependencies are autodetected and run appropriately
  hof/flow provides many built in task types
  you can reuse, combine, and share as CUE modules, packages, and values
`

func init() {

	FlowCmd.Flags().BoolVarP(&(flags.FlowFlags.List), "list", "l", false, "list available pipelines")
	FlowCmd.Flags().BoolVarP(&(flags.FlowFlags.Docs), "docs", "d", false, "print pipeline docs")
	FlowCmd.Flags().StringSliceVarP(&(flags.FlowFlags.Flow), "flow", "f", nil, "flow labels to match and run")
	FlowCmd.Flags().BoolVarP(&(flags.FlowFlags.Progress), "progress", "", false, "print task progress as it happens")
	FlowCmd.Flags().BoolVarP(&(flags.FlowFlags.Stats), "stats", "s", false, "Print final task statistics")
}

func FlowRun(globs []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var FlowCmd = &cobra.Command{

	Use: "flow [cue files...] [@flow/name...] [+key=value]",

	Aliases: []string{
		"f",
	},

	Short: "run CUE pipelines with the hof/flow DAG engine",

	Long: flowLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		var globs []string

		if 0 < len(args) {

			globs = args[0:]

		}

		err = FlowRun(globs)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := FlowCmd.HelpFunc()
	ousage := FlowCmd.UsageFunc()
	help := func(cmd *cobra.Command, args []string) {
		if extra(cmd) {
			return
		}
		ohelp(cmd, args)
	}
	usage := func(cmd *cobra.Command) error {
		if extra(cmd) {
			return nil
		}
		return ousage(cmd)
	}

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	FlowCmd.SetHelpFunc(thelp)
	FlowCmd.SetUsageFunc(tusage)

}
