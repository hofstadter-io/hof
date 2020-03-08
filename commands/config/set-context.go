package config

import (
	"fmt"
	"os"

	"github.com/hofstadter-io/hof/pkg/config"
	"github.com/spf13/cobra"
)

var SetContextLong = `Sets the context values and makes it the default`

var SetContextCmd = &cobra.Command{

	Use: "set-context <context-name> <username> <apikey> <host>",

	Short: "Set a context configuration",

	Long: SetContextLong,

	Run: func(cmd *cobra.Command, args []string) {

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'context'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var context string
		if 0 < len(args) {
			context = args[0]
		}

		if 1 >= len(args) {
			fmt.Println("missing required argument: 'account'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var account string
		if 1 < len(args) {
			account = args[1]
		}

		if 2 >= len(args) {
			fmt.Println("missing required argument: 'apikey'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var apikey string
		if 2 < len(args) {
			apikey = args[2]
		}

		var host string
		host = "https://studios.studios.live.hofstadter.io"
		if 3 < len(args) {
			host = args[3]
		}

		/*
			fmt.Println("hof config set-context:",
				context,
				account,
				apikey,
				host,
			)
		*/

		config.SetContext(context, account, apikey, host)
	},
}
