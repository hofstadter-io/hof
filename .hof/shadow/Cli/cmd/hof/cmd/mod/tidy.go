package cmdmod

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/lib/mod"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var tidyLong = `add missinad and remove unused modules`

func TidyRun(args []string) (err error) {

	err = mod.ProcessLangs("tidy", args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

var TidyCmd = &cobra.Command{

	Use: "tidy [langs...]",

	Short: "add missinad and remove unused modules",

	Long: tidyLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = TidyRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {

	help := TidyCmd.HelpFunc()
	usage := TidyCmd.UsageFunc()

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
	TidyCmd.SetHelpFunc(thelp)
	TidyCmd.SetUsageFunc(tusage)

}
