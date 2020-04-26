package cmd

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/studios"
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
	StudiosCmd.AddCommand(cmdstudios.AppCmd)
	StudiosCmd.AddCommand(cmdstudios.DatabaseCmd)
	StudiosCmd.AddCommand(cmdstudios.ContainerCmd)
	StudiosCmd.AddCommand(cmdstudios.FunctionCmd)
	StudiosCmd.AddCommand(cmdstudios.ConfigCmd)
	StudiosCmd.AddCommand(cmdstudios.SecretCmd)
}
