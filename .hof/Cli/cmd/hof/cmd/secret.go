package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/secret"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var secretLong = `secret subcommands`

var SecretCmd = &cobra.Command{

	Use: "secret",

	Short: "secret subcommands",

	Long: secretLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},
}

func init() {

	help := SecretCmd.HelpFunc()
	usage := SecretCmd.UsageFunc()

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
	SecretCmd.SetHelpFunc(thelp)
	SecretCmd.SetUsageFunc(tusage)

	SecretCmd.AddCommand(cmdsecret.CreateCmd)
	SecretCmd.AddCommand(cmdsecret.ListCmd)
	SecretCmd.AddCommand(cmdsecret.GetCmd)
	SecretCmd.AddCommand(cmdsecret.SetCmd)
	SecretCmd.AddCommand(cmdsecret.UseCmd)

}
