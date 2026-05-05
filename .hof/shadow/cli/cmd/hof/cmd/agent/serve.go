package cmdagent

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	libagentcmd "github.com/hofstadter-io/hof/lib/agent/cmd"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var serveLong = `serve agentic components`

func ServeRun(args []string) (err error) {

	err = libagentcmd.Serve(args, flags.RootPflags, flags.AgentPflags)

	return err
}

var ServeCmd = &cobra.Command{

	Use: "serve [...target] [% ...cue]",

	Short: "serve agentic components",

	Long: serveLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = ServeRun(args)
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

	ohelp := ServeCmd.HelpFunc()
	ousage := ServeCmd.UsageFunc()

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
	ServeCmd.SetHelpFunc(thelp)
	ServeCmd.SetUsageFunc(tusage)

}
