package studios

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios/container"
)

var containerLong = `Work with Hofstadter Studios containers`

var ContainerCmd = &cobra.Command{

	Use: "container",

	Aliases: []string{
		"cont",
	},

	Short: "Work with Hofstadter Studios containers",

	Long: containerLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios container")

	},
}

func init() {
	ContainerCmd.AddCommand(container.CallCmd)
	ContainerCmd.AddCommand(container.ListCmd)
	ContainerCmd.AddCommand(container.GetCmd)
	ContainerCmd.AddCommand(container.CreateCmd)
	ContainerCmd.AddCommand(container.UpdateCmd)
	ContainerCmd.AddCommand(container.DeployCmd)
	ContainerCmd.AddCommand(container.StatusCmd)
	ContainerCmd.AddCommand(container.PushCmd)
	ContainerCmd.AddCommand(container.PullCmd)
	ContainerCmd.AddCommand(container.ResetCmd)
	ContainerCmd.AddCommand(container.ShutdownCmd)
	ContainerCmd.AddCommand(container.DeleteCmd)
}
