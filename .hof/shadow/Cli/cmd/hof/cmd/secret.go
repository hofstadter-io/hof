package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/secret"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var secretLong = `manage local secrets`

var SecretCmd = &cobra.Command{

	Use: "secret",

	Short: "manage local secrets",

	Long: secretLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

	},
}

func init() {

	help := SecretCmd.HelpFunc()
	usage := SecretCmd.UsageFunc()

	thelp := func(cmd *cobra.Command, args []string) {
		ga.SendCommandPath(cmd.CommandPath() + " help")
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		ga.SendCommandPath(cmd.CommandPath() + " usage")
		return usage(cmd)
	}
	SecretCmd.SetHelpFunc(thelp)
	SecretCmd.SetUsageFunc(tusage)

	SecretCmd.AddCommand(cmdsecret.GetCmd)
	SecretCmd.AddCommand(cmdsecret.SetCmd)
	SecretCmd.AddCommand(cmdsecret.UseCmd)

}
