package cmdenv

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	libenvcmd "github.com/hofstadter-io/hof/lib/env/cmd"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var runLong = `run target point in an environment`

func init() {

	flags.SetupEnv__RunFlags(RunCmd.Flags(), &(flags.Env__RunFlags))

}

func RunRun(args []string) (err error) {

	err = libenvcmd.Run(args, flags.RootPflags, flags.EnvPflags, flags.Env__RunFlags)

	return err
}

var RunCmd = &cobra.Command{

	Use: "run <target> [% [...cue]]",

	Short: "run target point in an environment",

	Long: runLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

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
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := RunCmd.HelpFunc()
	ousage := RunCmd.UsageFunc()

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
	RunCmd.SetHelpFunc(thelp)
	RunCmd.SetUsageFunc(tusage)

}
