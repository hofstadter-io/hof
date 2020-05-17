package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/config"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var configLong = `configuration subcommands`

var ConfigCmd = &cobra.Command{

	Use: "config",

	Short: "configuration subcommands",

	Long: configLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},
}

func init() {

	help := ConfigCmd.HelpFunc()
	usage := ConfigCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		return usage(cmd)
	}
	ConfigCmd.SetHelpFunc(thelp)
	ConfigCmd.SetUsageFunc(tusage)

	ConfigCmd.AddCommand(cmdconfig.CreateCmd)
	ConfigCmd.AddCommand(cmdconfig.ListCmd)
	ConfigCmd.AddCommand(cmdconfig.GetCmd)
	ConfigCmd.AddCommand(cmdconfig.SetCmd)
	ConfigCmd.AddCommand(cmdconfig.UseCmd)

}
