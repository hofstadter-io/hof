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
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := SecretCmd.HelpFunc()
	ousage := SecretCmd.UsageFunc()
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
	SecretCmd.SetHelpFunc(thelp)
	SecretCmd.SetUsageFunc(tusage)

	SecretCmd.AddCommand(cmdsecret.GetCmd)
	SecretCmd.AddCommand(cmdsecret.SetCmd)
	SecretCmd.AddCommand(cmdsecret.UseCmd)

}
