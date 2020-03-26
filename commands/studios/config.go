package studios

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/commands/studios/config"
)

var configLong = `Work with Hofstadter Studios configs`

var ConfigCmd = &cobra.Command{

	Use: "config",

	Aliases: []string{
		"cfg",
	},

	Short: "Work with Hofstadter Studios configs",

	Long: configLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		// Default body

		fmt.Println("hof studios config")

	},
}

func init() {
	ConfigCmd.AddCommand(config.ListCmd)
	ConfigCmd.AddCommand(config.GetCmd)
	ConfigCmd.AddCommand(config.CreateCmd)
	ConfigCmd.AddCommand(config.UpdateCmd)
	ConfigCmd.AddCommand(config.DeleteCmd)
}
