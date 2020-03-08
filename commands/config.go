package commands

import (

	// custom imports

	// infered imports

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/config"
)

// Tool:   hof
// Name:   config
// Usage:  config <cmd>
// Parent: hof

var ConfigLong = `Configure the Hof CLI`

var ConfigCmd = &cobra.Command{

	Use: "config <cmd>",

	Short: "Configure the Hof CLI",

	Long: ConfigLong,
}

func init() {
	RootCmd.AddCommand(ConfigCmd)
}

func init() {
	// add sub-commands to this command when present

	ConfigCmd.AddCommand(config.TestCmd)
	ConfigCmd.AddCommand(config.UseContextCmd)
	ConfigCmd.AddCommand(config.SetContextCmd)
	ConfigCmd.AddCommand(config.GetCmd)
	ConfigCmd.AddCommand(config.SetCmd)
}
