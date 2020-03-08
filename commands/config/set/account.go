package set

import (
	"fmt"

	// custom imports

	// infered imports
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/hofstadter-io/hof/lib/config"
)

// Tool:   hof
// Name:   account
// Usage:  account <name>
// Parent: set

var AccountLong = `Set your account ID`

var AccountCmd = &cobra.Command{

	Use: "account <name>",

	Short: "Set your account ID",

	Long: AccountLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In accountCmd", "args", args)
		// Argument Parsing
		// [0]name:   account
		//     help:
		//     req'd:  true
		if 0 >= len(args) {
			fmt.Println("missing required argument: 'account'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var account string

		if 0 < len(args) {

			account = args[0]
		}

		/*
		fmt.Println("hof config set account:",
			account,
		)
		*/

		context := viper.GetString("context")
		config.SetAccount(context, account)
	},
}

func init() {
	// add sub-commands to this command when present

}
