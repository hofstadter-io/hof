package set

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hofstadter-io/hof/pkg/config"
)

var HostLong = `Set your host server`

var HostCmd = &cobra.Command{

	Use: "host <name>",

	Short: "Set your host server",

	Long: HostLong,

	Run: func(cmd *cobra.Command, args []string) {

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'host'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var host string
		if 0 < len(args) {
			host = args[0]
		}

		/*
			fmt.Println("hof config set host:",
				host,
			)
		*/

		context := viper.GetString("context")
		config.SetHost(context, host)
	},
}
