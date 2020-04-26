package cmdstudios

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/studios/function"
)

var functionLong = `Work with Hofstadter Studios functions`

func FunctionRun(args []string) (err error) {

	return err
}

var FunctionCmd = &cobra.Command{

	Use: "function",

	Aliases: []string{
		"cont",
	},

	Short: "Work with Hofstadter Studios functions",

	Long: functionLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = FunctionRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	FunctionCmd.AddCommand(cmdfunction.CallCmd)
	FunctionCmd.AddCommand(cmdfunction.ListCmd)
	FunctionCmd.AddCommand(cmdfunction.GetCmd)
	FunctionCmd.AddCommand(cmdfunction.CreateCmd)
	FunctionCmd.AddCommand(cmdfunction.UpdateCmd)
	FunctionCmd.AddCommand(cmdfunction.DeployCmd)
	FunctionCmd.AddCommand(cmdfunction.StatusCmd)
	FunctionCmd.AddCommand(cmdfunction.PushCmd)
	FunctionCmd.AddCommand(cmdfunction.PullCmd)
	FunctionCmd.AddCommand(cmdfunction.ResetCmd)
	FunctionCmd.AddCommand(cmdfunction.ShutdownCmd)
	FunctionCmd.AddCommand(cmdfunction.DeleteCmd)
}
