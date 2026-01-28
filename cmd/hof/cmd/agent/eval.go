package cmdagent

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	libagentcmd "github.com/hofstadter-io/hof/lib/agent/cmd"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var evalLong = `eval agents over parameters`

func EvalRun(args []string) (err error) {

	err = libagentcmd.Eval(args, flags.RootPflags, flags.AgentPflags)

	return err
}

var EvalCmd = &cobra.Command{

	Use: "eval [...target] [% ...cue]",

	Short: "eval agents over parameters",

	Long: evalLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = EvalRun(args)
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

	ohelp := EvalCmd.HelpFunc()
	ousage := EvalCmd.UsageFunc()

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
	EvalCmd.SetHelpFunc(thelp)
	EvalCmd.SetUsageFunc(tusage)

}
