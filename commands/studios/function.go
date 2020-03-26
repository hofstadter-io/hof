package studios

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios/function"
)

var functionLong = `Work with Hofstadter Studios functions`

var FunctionCmd = &cobra.Command{

	Use: "function",

	Aliases: []string{
		"cont",
	},

	Short: "Work with Hofstadter Studios functions",

	Long: functionLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios function")

	},
}

func init() {
	FunctionCmd.AddCommand(function.CallCmd)
	FunctionCmd.AddCommand(function.ListCmd)
	FunctionCmd.AddCommand(function.GetCmd)
	FunctionCmd.AddCommand(function.CreateCmd)
	FunctionCmd.AddCommand(function.UpdateCmd)
	FunctionCmd.AddCommand(function.DeployCmd)
	FunctionCmd.AddCommand(function.StatusCmd)
	FunctionCmd.AddCommand(function.PushCmd)
	FunctionCmd.AddCommand(function.PullCmd)
	FunctionCmd.AddCommand(function.ResetCmd)
	FunctionCmd.AddCommand(function.ShutdownCmd)
	FunctionCmd.AddCommand(function.DeleteCmd)
}
