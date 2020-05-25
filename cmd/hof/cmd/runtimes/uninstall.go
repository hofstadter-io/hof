package cmdruntimes

import (
	"fmt"
	"os"

	"strings"

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

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

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

	help := UninstallCmd.HelpFunc()
	usage := UninstallCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/usage", "<omit>", 0)
		return usage(cmd)
	}
	UninstallCmd.SetHelpFunc(thelp)
	UninstallCmd.SetUsageFunc(tusage)

}
