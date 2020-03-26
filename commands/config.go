package commands

import (

	// hello... something might need to go here

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/config"
)

var configLong = `configuration subcommands`

var ConfigCmd = &cobra.Command{

	Use: "config",

	Short: "configuration subcommands",

	Long: configLong,
}

func init() {
	ConfigCmd.AddCommand(config.TestCmd)
	ConfigCmd.AddCommand(config.ListCmd)
	ConfigCmd.AddCommand(config.GetCmd)
	ConfigCmd.AddCommand(config.SetCmd)
	ConfigCmd.AddCommand(config.UseCmd)
}
