package studios

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios/functions"
)

var functionsLong = `Work with Hofstadter Studios functions`

var FunctionsCmd = &cobra.Command{

	Use: "functions",

	Aliases: []string{
		"cont",
	},

	Short: "Work with Hofstadter Studios functions",

	Long: functionsLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios functions")

	},
}

func init() {
	FunctionsCmd.AddCommand(functions.CallCmd)
	FunctionsCmd.AddCommand(functions.ListCmd)
	FunctionsCmd.AddCommand(functions.GetCmd)
	FunctionsCmd.AddCommand(functions.CreateCmd)
	FunctionsCmd.AddCommand(functions.UpdateCmd)
	FunctionsCmd.AddCommand(functions.DeployCmd)
	FunctionsCmd.AddCommand(functions.StatusCmd)
	FunctionsCmd.AddCommand(functions.PushCmd)
	FunctionsCmd.AddCommand(functions.PullCmd)
	FunctionsCmd.AddCommand(functions.ResetCmd)
	FunctionsCmd.AddCommand(functions.ShutdownCmd)
	FunctionsCmd.AddCommand(functions.DeleteCmd)
}
