package cmdstudios

import (
	"fmt"
	"os"

	"strings"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/studios/function"

	"github.com/hofstadter-io/hof/cmd/hof/ga"
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

	PreRun: func(cmd *cobra.Command, args []string) {

		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		ga.SendGaEvent(c, strings.Join(args, "/"), 0)

	},

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
	hf := FunctionCmd.HelpFunc()
	f := func(cmd *cobra.Command, args []string) {
		cs := strings.Fields(cmd.CommandPath())
		c := strings.Join(cs[1:], "/")
		as := strings.Join(args, "/")
		ga.SendGaEvent(c+"/help", as, 0)
		hf(cmd, args)
	}
	FunctionCmd.SetHelpFunc(f)
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
