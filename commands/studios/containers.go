package studios

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios/containers"
)

var containersLong = `Work with Hofstadter Studios containers`

var ContainersCmd = &cobra.Command{

	Use: "containers",

	Aliases: []string{
		"cont",
	},

	Short: "Work with Hofstadter Studios containers",

	Long: containersLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios containers")

	},
}

func init() {
	ContainersCmd.AddCommand(containers.CallCmd)
	ContainersCmd.AddCommand(containers.ListCmd)
	ContainersCmd.AddCommand(containers.GetCmd)
	ContainersCmd.AddCommand(containers.CreateCmd)
	ContainersCmd.AddCommand(containers.UpdateCmd)
	ContainersCmd.AddCommand(containers.DeployCmd)
	ContainersCmd.AddCommand(containers.StatusCmd)
	ContainersCmd.AddCommand(containers.PushCmd)
	ContainersCmd.AddCommand(containers.PullCmd)
	ContainersCmd.AddCommand(containers.ResetCmd)
	ContainersCmd.AddCommand(containers.ShutdownCmd)
	ContainersCmd.AddCommand(containers.DeleteCmd)
}
