package studios

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios/app"
)

var appLong = `Work with Hofstadter Studios apps`

var AppCmd = &cobra.Command{

	Use: "app",

	Short: "Work with Hofstadter Studios apps",

	Long: appLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios app")

	},
}

func init() {
	AppCmd.AddCommand(app.ListCmd)
	AppCmd.AddCommand(app.GetCmd)
	AppCmd.AddCommand(app.CreateCmd)
	AppCmd.AddCommand(app.UpdateCmd)
	AppCmd.AddCommand(app.DeployCmd)
	AppCmd.AddCommand(app.StatusCmd)
	AppCmd.AddCommand(app.PushCmd)
	AppCmd.AddCommand(app.PullCmd)
	AppCmd.AddCommand(app.ResetCmd)
	AppCmd.AddCommand(app.ShutdownCmd)
	AppCmd.AddCommand(app.DeleteCmd)
}
