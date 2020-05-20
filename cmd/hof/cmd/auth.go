package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/auth"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
)

var authLong = `authentication subcommands`

var AuthCmd = &cobra.Command{

	Use: "auth",

	Short: "authentication subcommands",

	Long: authLong,

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, "<omit>", 0)

	},
}

func init() {

	help := AuthCmd.HelpFunc()
	usage := AuthCmd.UsageFunc()

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
	AuthCmd.SetHelpFunc(thelp)
	AuthCmd.SetUsageFunc(tusage)

	AuthCmd.AddCommand(cmdauth.LoginCmd)
	AuthCmd.AddCommand(cmdauth.LogoutCmd)
	AuthCmd.AddCommand(cmdauth.ListCmd)
	AuthCmd.AddCommand(cmdauth.TestCmd)

}
