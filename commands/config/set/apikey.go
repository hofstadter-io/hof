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
// Name:   apikey
// Usage:  apikey <key>
// Parent: set

var ApikeyLong = `Set your API Key`

var ApikeyCmd = &cobra.Command{

	Use: "apikey <key>",

	Short: "Set your API Key",

	Long: ApikeyLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In apikeyCmd", "args", args)
		// Argument Parsing
		// [0]name:   key
		//     help:
		//     req'd:  true
		if 0 >= len(args) {
			fmt.Println("missing required argument: 'key'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var key string

		if 0 < len(args) {

			key = args[0]
		}

		/*
		fmt.Println("hof config set apikey:",
			key,
		)
		*/

		context := viper.GetString("context")
		config.SetAPIKey(context, key)
	},
}

func init() {
	// add sub-commands to this command when present

}
