package cmdruntimes

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/runtimes"
)

var installLong = `install a runtime`

func InstallRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = runtimes.RunInstallFromArgs(args)

	return err
}

var InstallCmd = &cobra.Command{

	Use: "install",

	Aliases: []string{
		"I",
	},

	Short: "install a runtime",

	Long: installLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = InstallRun(args)
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

	ohelp := InstallCmd.HelpFunc()
	ousage := InstallCmd.UsageFunc()
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
	InstallCmd.SetHelpFunc(thelp)
	InstallCmd.SetUsageFunc(tusage)

}