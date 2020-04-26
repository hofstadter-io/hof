package cmdstudios

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/studios/app"
)

var appLong = `Work with Hofstadter Studios apps`

func AppRun(args []string) (err error) {

	return err
}

var AppCmd = &cobra.Command{

	Use: "app",

	Short: "Work with Hofstadter Studios apps",

	Long: appLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = AppRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	AppCmd.AddCommand(cmdapp.ListCmd)
	AppCmd.AddCommand(cmdapp.GetCmd)
	AppCmd.AddCommand(cmdapp.CreateCmd)
	AppCmd.AddCommand(cmdapp.UpdateCmd)
	AppCmd.AddCommand(cmdapp.DeployCmd)
	AppCmd.AddCommand(cmdapp.StatusCmd)
	AppCmd.AddCommand(cmdapp.PushCmd)
	AppCmd.AddCommand(cmdapp.PullCmd)
	AppCmd.AddCommand(cmdapp.ResetCmd)
	AppCmd.AddCommand(cmdapp.ShutdownCmd)
	AppCmd.AddCommand(cmdapp.DeleteCmd)
}
