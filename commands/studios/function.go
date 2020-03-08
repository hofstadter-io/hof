package studios

import (
	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios/function"
)

var FunctionLong = `Work with your Studios Functions`

var FunctionCmd = &cobra.Command{

	Use: "function",

	Aliases: []string{
		"functions",
		"funcs",
		"func",
		"fns",
		"fn",
	},

	Short: "Work with your Studios Functions",

	Long: FunctionLong,
}

func init() {
	// add sub-commands
	FunctionCmd.AddCommand(function.VersionsCmd)
	FunctionCmd.AddCommand(function.StatusCmd)
	FunctionCmd.AddCommand(function.LogsCmd)
	FunctionCmd.AddCommand(function.ListCmd)
	FunctionCmd.AddCommand(function.DeployCmd)
	FunctionCmd.AddCommand(function.CreateCmd)
	FunctionCmd.AddCommand(function.DeleteCmd)
	FunctionCmd.AddCommand(function.DeployCmd)
	FunctionCmd.AddCommand(function.ShutdownCmd)
	FunctionCmd.AddCommand(function.CallCmd)
	FunctionCmd.AddCommand(function.PullCmd)
	FunctionCmd.AddCommand(function.PushCmd)
}
