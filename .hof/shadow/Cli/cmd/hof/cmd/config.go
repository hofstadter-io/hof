package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/config"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var configLong = `manage local configurations`

var ConfigCmd = &cobra.Command{

	Use: "config",

	Short: "manage local configurations",

	Long: configLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},
}

func init() {

	help := ConfigCmd.HelpFunc()
	usage := ConfigCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	ConfigCmd.SetHelpFunc(thelp)
	ConfigCmd.SetUsageFunc(tusage)

	ConfigCmd.AddCommand(cmdconfig.GetCmd)
	ConfigCmd.AddCommand(cmdconfig.SetCmd)
	ConfigCmd.AddCommand(cmdconfig.UseCmd)

}
