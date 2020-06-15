package cmdruntimes

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/runtimes"
)

var uninstallLong = `uninstall a runtime`

func UninstallRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = runtimes.RunUninstallFromArgs(args)

	return err
}

var UninstallCmd = &cobra.Command{

	Use: "uninstall",

	Short: "uninstall a runtime",

	Long: uninstallLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = UninstallRun(args)
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

	ohelp := UninstallCmd.HelpFunc()
	ousage := UninstallCmd.UsageFunc()
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
	UninstallCmd.SetHelpFunc(thelp)
	UninstallCmd.SetUsageFunc(tusage)

}