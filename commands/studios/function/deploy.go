package function

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hofstadter-io/hof/pkg/studios/function"
)

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

		// fmt.Println("hof function deploy:")

		err := function.Deploy(DeployPushFlag, DeployMemoryFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
