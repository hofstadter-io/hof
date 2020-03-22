package commands

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios"
)

var StudiosLong = `Studios subcommands for the Hof CLI`

var StudiosCmd = &cobra.Command{

	Use: "studios <cmd>",

	Short: "Studios subcommands for the Hof CLI",

	Long: StudiosLong,
}

func init() {
	RootCmd.AddCommand(StudiosCmd)
}

func init() {
	// add sub-commands to this command when present

	StudiosCmd.AddCommand(studios.AppCmd)
	StudiosCmd.AddCommand(studios.ContainerCmd)
	StudiosCmd.AddCommand(studios.DatabaseCmd)
	StudiosCmd.AddCommand(studios.FunctionCmd)
	StudiosCmd.AddCommand(studios.SecretCmd)
}
