package set

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hofstadter-io/hof/pkg/config"
)

var ApikeyLong = `Set your API Key`

var ApikeyCmd = &cobra.Command{

	Use: "apikey <key>",

	Short: "Set your API Key",

	Long: ApikeyLong,

	Run: func(cmd *cobra.Command, args []string) {

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
