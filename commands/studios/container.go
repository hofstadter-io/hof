package studios

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios/container"
)

var ContainerLong = `Work with your Studios Container Run`

var ContainerCmd = &cobra.Command{

	Use: "container",

	Aliases: []string{
		"containers",
		"crun",
	},

	Short: "Work with your Studios Container Run",

	Long: ContainerLong,
}

func init() {
	// add sub-commands
	ContainerCmd.AddCommand(container.StatusCmd)
	ContainerCmd.AddCommand(container.LogsCmd)
	ContainerCmd.AddCommand(container.ListCmd)
	ContainerCmd.AddCommand(container.CreateCmd)
	ContainerCmd.AddCommand(container.DeleteCmd)
	ContainerCmd.AddCommand(container.ShutdownCmd)
	ContainerCmd.AddCommand(container.CallCmd)
	ContainerCmd.AddCommand(container.PullCmd)
	ContainerCmd.AddCommand(container.PushCmd)
	ContainerCmd.AddCommand(container.DeployCmd)
}
