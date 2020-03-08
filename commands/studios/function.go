package commands

import (

	// custom imports

	// infered imports

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/function"
)

// Tool:   hof
// Name:   function
// Usage:  function
// Parent: hof

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
	RootCmd.AddCommand(FunctionCmd)
}

func init() {
	// add sub-commands to this command when present

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
