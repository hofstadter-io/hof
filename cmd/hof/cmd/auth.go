package cmd

import (
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

		ga.SendCommandPath(cmd.CommandPath())

	},
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := AuthCmd.HelpFunc()
	ousage := AuthCmd.UsageFunc()
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
	AuthCmd.SetHelpFunc(thelp)
	AuthCmd.SetUsageFunc(tusage)

	AuthCmd.AddCommand(cmdauth.LoginCmd)
	AuthCmd.AddCommand(cmdauth.LogoutCmd)
	AuthCmd.AddCommand(cmdauth.ListCmd)
	AuthCmd.AddCommand(cmdauth.TestCmd)

}
