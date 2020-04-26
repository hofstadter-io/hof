package cmdstudios

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/cmd/hof/cmd/studios/config"
)

var configLong = `Work with Hofstadter Studios configs`

func ConfigRun(args []string) (err error) {

	return err
}

var ConfigCmd = &cobra.Command{

	Use: "config",

	Aliases: []string{
		"cfg",
	},

	Short: "Work with Hofstadter Studios configs",

	Long: configLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// Argument Parsing

		err = ConfigRun(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(cmdconfig.ListCmd)
	ConfigCmd.AddCommand(cmdconfig.GetCmd)
	ConfigCmd.AddCommand(cmdconfig.CreateCmd)
	ConfigCmd.AddCommand(cmdconfig.UpdateCmd)
	ConfigCmd.AddCommand(cmdconfig.DeleteCmd)
}
