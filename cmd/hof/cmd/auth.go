package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/auth"
)

var authLong = `authentication subcommands`

var AuthCmd = &cobra.Command{

	Use: "auth",

	Short: "authentication subcommands",

	Long: authLong,
}

func init() {
	AuthCmd.AddCommand(cmdauth.LoginCmd)
}
