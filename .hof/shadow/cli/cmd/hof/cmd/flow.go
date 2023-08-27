package cmd

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/flow"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var flowLong = `run workflows and tasks powered by CUE`

func init() {

	FlowCmd.PersistentFlags().StringSliceVarP(&(flags.FlowPflags.Flow), "flow", "F", nil, "flow labels to match and run")
	FlowCmd.PersistentFlags().StringVarP(&(flags.FlowPflags.Bulk), "bulk", "B", "", "exprs for inputs to run workflow in bulk mode")
	FlowCmd.PersistentFlags().IntVarP(&(flags.FlowPflags.Parallel), "parallel", "P", 1, "bulk processing parallelism")
	FlowCmd.PersistentFlags().BoolVarP(&(flags.FlowPflags.Progress), "progress", "", false, "print task progress as it happens")
}

func FlowRun(args []string) (err error) {

	// you can safely comment this print out
	fmt.Println("not implemented")

	return err
}

var FlowCmd = &cobra.Command{

	Use: "flow [cue files...] [@flow/name...] [+key=value]",

	Aliases: []string{
		"f",
	},

	Short: "run workflows and tasks powered by CUE",

	Long: flowLong,

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		glob := toComplete + "*"
		matches, _ := filepath.Glob(glob)
		return matches, cobra.ShellCompDirectiveDefault
	},

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = FlowRun(args)
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

		ga.SendCommandPath(cmd.CommandPath() + " help")

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
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		return usage(cmd)
	}
	FlowCmd.SetHelpFunc(thelp)
	FlowCmd.SetUsageFunc(tusage)

	FlowCmd.AddCommand(cmdflow.ListCmd)

}
