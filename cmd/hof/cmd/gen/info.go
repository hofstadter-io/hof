package cmdgen

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/flags"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/gen/cmd"
)

var infoLong = `print details for about generators`

func init() {

	InfoCmd.Flags().StringSliceVarP(&(flags.Gen__InfoFlags.Expression), "expr", "e", nil, "CUE paths to select outputs, depending on the command")
}

func InfoRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = cmd.Info(args, flags.RootPflags, flags.GenFlags, flags.Gen__InfoFlags)

	return err
}

var InfoCmd = &cobra.Command{

	Use: "info",

	Short: "print details for about generators",

	Long: infoLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = InfoRun(args)
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
	InfoCmd.SetHelpFunc(thelp)
	InfoCmd.SetUsageFunc(tusage)

}
