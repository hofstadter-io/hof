package function

import (
	"fmt"
	"os"

	// custom imports

	// infered imports

	"github.com/hofstadter-io/hof/lib/fns"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   deploy
// Usage:  deploy
// Parent: function

var DeployLong = `Deploy the function <name> from the current directory`

var (
	DeployPushFlag   bool
	DeployMemoryFlag int
)

func init() {
	DeployCmd.Flags().BoolVarP(&DeployPushFlag, "push", "p", true, "push the latest function code with the deploy.")
	viper.BindPFlag("push", DeployCmd.Flags().Lookup("push"))

	DeployCmd.Flags().IntVarP(&DeployMemoryFlag, "memory", "m", 0, "set the memory for this service (in megabytes).")
	viper.BindPFlag("memory", DeployCmd.Flags().Lookup("memory"))

}

var DeployCmd = &cobra.Command{

	Use: "deploy",

	Short: "Deploys the function <name>",

	Long: DeployLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In deployCmd", "args", args)
		// Argument Parsing

		// fmt.Println("hof function deploy:")

		err := fns.Deploy(DeployPushFlag, DeployMemoryFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	// add sub-commands to this command when present

}
