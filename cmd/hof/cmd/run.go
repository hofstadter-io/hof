package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/ops"
	"github.com/hofstadter-io/hof/cmd/hof/flags"
)

var runLong = `HofLineScript (HLS) run polyglot command and scripts seamlessly across runtimes

can accept cue & flags or just a .hls file
`

func init() {

	RunCmd.Flags().BoolVarP(&(flags.RunFlags.List), "list", "", false, "list matching scripts that would run")
	RunCmd.Flags().BoolVarP(&(flags.RunFlags.Info), "info", "", false, "view detailed info for matching scripts")
	RunCmd.Flags().StringSliceVarP(&(flags.RunFlags.Suite), "suite", "s", nil, "<name>: _ @run(suite)'s to run")
	RunCmd.Flags().StringSliceVarP(&(flags.RunFlags.Runner), "runner", "r", nil, "<name>: _ @run(script)'s to run")
	RunCmd.Flags().StringSliceVarP(&(flags.RunFlags.Environment), "env", "e", nil, "exrta environment variables for scripts")
	RunCmd.Flags().StringSliceVarP(&(flags.RunFlags.Data), "data", "d", nil, "exrta data to include in the scripts context")
	RunCmd.Flags().StringVarP(&(flags.RunFlags.Workdir), "workdir", "w", "", "working directory")
}

func RunRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = ops.RunRunFromArgs(args)

	return err
}

var RunCmd = &cobra.Command{

	Use: "run",

	Aliases: []string{
		"r",
	},

	Short: "Hof Line Script (HLS) is a successor to bash and python based scripting",

	Long: runLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = RunRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := RunCmd.HelpFunc()
	usage := RunCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	RunCmd.SetHelpFunc(thelp)
	RunCmd.SetUsageFunc(tusage)

}