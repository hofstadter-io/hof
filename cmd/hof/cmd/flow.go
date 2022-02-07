package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/flow"
)

var flowLong = `run file(s) through the hof/flow DAG engine`

func init() {

	FlowCmd.Flags().BoolVarP(&(flags.FlowFlags.List), "list", "l", false, "list available pipelines")
	FlowCmd.Flags().BoolVarP(&(flags.FlowFlags.Docs), "docs", "d", false, "print pipeline docs")
	FlowCmd.Flags().StringSliceVarP(&(flags.FlowFlags.Flow), "flow", "f", nil, "flow labels to match and run")
	FlowCmd.Flags().StringSliceVarP(&(flags.FlowFlags.Tags), "tags", "t", nil, "data tags to inject before run")
	FlowCmd.Flags().BoolVarP(&(flags.FlowFlags.DebugTasks), "debug-tasks", "", false, "print debugging info about tasks")
}

func FlowRun(entrypoints []string) (err error) {
  if len(entrypoints) == 0 {
    entrypoints = []string{"./"}
  }
	return flow.Run(entrypoints, &flags.RootPflags, &flags.FlowFlags)
}

var FlowCmd = &cobra.Command{

	Use: "flow [cue files...]",

	Aliases: []string{
		"f",
	},

	Short: "run file(s) through the hof/flow DAG engine",

	Long: flowLong,

	PreRun: func(cmd *cobra.Command, args []string) {

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

	FlowCmd.SetHelpFunc(help)
	FlowCmd.SetUsageFunc(usage)

}
