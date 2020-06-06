package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/runtimes"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var runtimesLong = `work with runtimes (go, js, py, bash, custom)`

var RuntimesCmd = &cobra.Command{

	Use: "runtimes",

	Short: "work with runtimes (go, js, py, bash, custom)",

	Long: runtimesLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "", 0)

	},
}

func init() {

	help := RuntimesCmd.HelpFunc()
	usage := RuntimesCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/usage", "", 0)
		return usage(cmd)
	}
	RuntimesCmd.SetHelpFunc(thelp)
	RuntimesCmd.SetUsageFunc(tusage)

	RuntimesCmd.AddCommand(cmdruntimes.InfoCmd)
	RuntimesCmd.AddCommand(cmdruntimes.CreateCmd)
	RuntimesCmd.AddCommand(cmdruntimes.GetCmd)
	RuntimesCmd.AddCommand(cmdruntimes.SetCmd)
	RuntimesCmd.AddCommand(cmdruntimes.EditCmd)
	RuntimesCmd.AddCommand(cmdruntimes.DeleteCmd)
	RuntimesCmd.AddCommand(cmdruntimes.InstallCmd)
	RuntimesCmd.AddCommand(cmdruntimes.UninstallCmd)

}
