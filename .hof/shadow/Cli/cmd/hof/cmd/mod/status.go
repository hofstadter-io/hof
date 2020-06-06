package cmdmod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/mod"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var statusLong = `print module dependencies status`

func StatusRun(args []string) (err error) {

	err = mod.ProcessLangs("status", args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

var StatusCmd = &cobra.Command{

	Use: "status",

	Short: "print module dependencies status",

	Long: statusLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = StatusRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := StatusCmd.HelpFunc()
	usage := StatusCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	StatusCmd.SetHelpFunc(thelp)
	StatusCmd.SetUsageFunc(tusage)

}
