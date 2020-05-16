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
	hf := SecretCmd.HelpFunc()
	f := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c+"/help", "<omit>", 0)
		hf(cmd, args)
	}
	SecretCmd.SetHelpFunc(f)
	SecretCmd.AddCommand(cmdsecret.CreateCmd)
	SecretCmd.AddCommand(cmdsecret.ListCmd)
	SecretCmd.AddCommand(cmdsecret.GetCmd)
	SecretCmd.AddCommand(cmdsecret.SetCmd)
	SecretCmd.AddCommand(cmdsecret.UseCmd)
}
