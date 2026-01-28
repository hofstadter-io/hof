package cmdchat

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	libagentcmd "github.com/hofstadter-io/hof/lib/agent/cmd"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var deleteLong = `delete sessions for given targets`

func DeleteRun(args []string) (err error) {

	err = libagentcmd.ChatList(args, flags.RootPflags, flags.AgentPflags, flags.Agent__ChatPflags)

	return err
}

var DeleteCmd = &cobra.Command{

	Use: "delete [...target] [% ...cue]",

	Short: "delete sessions for given targets",

	Long: deleteLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = DeleteRun(args)
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

	ohelp := DeleteCmd.HelpFunc()
	ousage := DeleteCmd.UsageFunc()

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
	DeleteCmd.SetHelpFunc(thelp)
	DeleteCmd.SetUsageFunc(tusage)

}
