package container

import (
	"fmt"
	"os"

	// custom imports

	// infered imports

	"github.com/hofstadter-io/hof/lib/crun"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   deploy
// Usage:  deploy [name]
// Parent: container

var DeployLong = `Deploys a container by name or from the current directory`

var (
	DeployPushFlag bool

	DeployMemoryFlag string

	DeployConcurrencyFlag int

	DeployTimeoutFlag string

	DeployEnvsFlag string
)

func init() {
	DeployCmd.Flags().BoolVarP(&DeployPushFlag, "push", "p", true, "push the latest code with the deploy.")
	viper.BindPFlag("push", DeployCmd.Flags().Lookup("push"))

	DeployCmd.Flags().StringVarP(&DeployMemoryFlag, "memory", "m", "128Mi", "sets the amount of memory.")
	viper.BindPFlag("memory", DeployCmd.Flags().Lookup("memory"))

	DeployCmd.Flags().IntVarP(&DeployConcurrencyFlag, "concurrency", "c", 80, "sets the number of concurrent requests.")
	viper.BindPFlag("concurrency", DeployCmd.Flags().Lookup("concurrency"))

	DeployCmd.Flags().StringVarP(&DeployTimeoutFlag, "timeout", "t", "60s", "sets the timeout for a request.")
	viper.BindPFlag("timeout", DeployCmd.Flags().Lookup("timeout"))

	DeployCmd.Flags().StringVarP(&DeployEnvsFlag, "envs", "e", "", "set env vars KEY=VALUE,[KEY=VALUE]...")
	viper.BindPFlag("envs", DeployCmd.Flags().Lookup("envs"))

}

var DeployCmd = &cobra.Command{

	Use: "deploy [name]",

	Short: "Deploys a container by name or from the current directory",

	Long: DeployLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In deployCmd", "args", args)
		// Argument Parsing
		// [0]name:   name
		//     help:
		//     req'd:

		var name string

		if 0 < len(args) {

			name = args[0]
		}

		/*
			fmt.Println("hof containers deploy:",
				name,
			)
		*/

		err := crun.Deploy(name, DeployPushFlag, DeployMemoryFlag, DeployConcurrencyFlag, DeployTimeoutFlag, DeployEnvsFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	// add sub-commands to this command when present

}
