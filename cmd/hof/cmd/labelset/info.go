package cmdlabelset

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/ga"

	"github.com/hofstadter-io/hof/lib/labels"
)

var infoLong = `print info about labelsets in your workspace or system`

func InfoRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = labels.RunInfoLabelsetFromArgs(args)

	return err
}

var InfoCmd = &cobra.Command{

	Use: "info",

	Aliases: []string{
		"i",
	},

	Short: "print info about labelsets in your workspace or system",

	Long: infoLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

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

	help := InfoCmd.HelpFunc()
	usage := InfoCmd.UsageFunc()

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
	InfoCmd.SetHelpFunc(thelp)
	InfoCmd.SetUsageFunc(tusage)

}
