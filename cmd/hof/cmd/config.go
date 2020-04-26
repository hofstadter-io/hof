package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/config"
)

var configLong = `configuration subcommands`

var ConfigCmd = &cobra.Command{

	Use: "config",

	Short: "configuration subcommands",

	Long: configLong,
}

func init() {
	ConfigCmd.AddCommand(cmdconfig.TestCmd)
	ConfigCmd.AddCommand(cmdconfig.ListCmd)
	ConfigCmd.AddCommand(cmdconfig.GetCmd)
	ConfigCmd.AddCommand(cmdconfig.SetCmd)
	ConfigCmd.AddCommand(cmdconfig.UseCmd)
}
