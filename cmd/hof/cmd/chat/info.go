package cmdchat

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"
	"github.com/hofstadter-io/hof/cmd/hof/ga"
	"github.com/hofstadter-io/hof/lib/chat/cmd"
)

var infoLong = `print details of a specific chat plugin`

func InfoRun(name string, entrypoints []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = cmd.Info(name, entrypoints, flags.RootPflags)

	return err
}

var InfoCmd = &cobra.Command{

	Use: "info [name]",

	Short: "print details of a specific chat plugin",

	Long: infoLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'name'")
			cmd.Usage()
			os.Exit(1)
		}

		var name string

		if 0 < len(args) {

			name = args[0]

		}

		var entrypoints []string

		if 1 < len(args) {

			entrypoints = args[1:]

		}

		err = InfoRun(name, entrypoints)
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

	ohelp := InfoCmd.HelpFunc()
	ousage := InfoCmd.UsageFunc()

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
	InfoCmd.SetHelpFunc(thelp)
	InfoCmd.SetUsageFunc(tusage)

}
