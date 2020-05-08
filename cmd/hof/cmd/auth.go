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
		ga.SendGaEvent(c, strings.Join(args, "/"), 0)

	},
}

func init() {
	hf := AuthCmd.HelpFunc()
	f := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		as := strings.Join(args, "/")
		ga.SendGaEvent(c+"/help", as, 0)
		hf(cmd, args)
	}
	AuthCmd.SetHelpFunc(f)
	AuthCmd.AddCommand(cmdauth.LoginCmd)
}
