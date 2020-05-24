package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/runtimes"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var runtimesLong = `work with runtimes`

var RuntimesCmd = &cobra.Command{

	Use: "runtimes",

	Short: "work with runtimes",

	Long: runtimesLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},
}

func init() {

	help := RuntimesCmd.HelpFunc()
	usage := RuntimesCmd.UsageFunc()

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
	RuntimesCmd.SetHelpFunc(thelp)
	RuntimesCmd.SetUsageFunc(tusage)

	RuntimesCmd.AddCommand(cmdruntimes.InfoCmd)
	RuntimesCmd.AddCommand(cmdruntimes.AddCmd)
	RuntimesCmd.AddCommand(cmdruntimes.GetCmd)
	RuntimesCmd.AddCommand(cmdruntimes.EditCmd)
	RuntimesCmd.AddCommand(cmdruntimes.RemoveCmd)

}
