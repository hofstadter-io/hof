package commands

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios"
)

var studiosLong = `  Hofstadter Studios makes it easy to develop and launch both
  hof-lang modules as well as pretty much any code or application`

var StudiosCmd = &cobra.Command{

	Use: "studios",

	Aliases: []string{
		"s",
	},

	Short: "commands for working with Hofstadter Studios",

	Long: studiosLong,
}

func init() {
	StudiosCmd.AddCommand(studios.AppCmd)
	StudiosCmd.AddCommand(studios.DatabaseCmd)
	StudiosCmd.AddCommand(studios.ContainerCmd)
	StudiosCmd.AddCommand(studios.FunctionCmd)
	StudiosCmd.AddCommand(studios.ConfigCmd)
	StudiosCmd.AddCommand(studios.SecretCmd)
}
