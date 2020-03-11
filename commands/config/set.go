package config

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/config/set"
)

var SetLong = `Get configuration values`

var SetCmd = &cobra.Command{

	Use: "set",

	Short: "Get configuration values",

	Long: SetLong,
}

func init() {
	// add sub-commands to this command when present

	SetCmd.AddCommand(set.ApikeyCmd)
	SetCmd.AddCommand(set.AccountCmd)
	SetCmd.AddCommand(set.HostCmd)
}
