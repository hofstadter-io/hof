package config

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

var useLong = `set the default configuration`

var UseCmd = &cobra.Command{

	Use: "use",

	Short: "set the default configuration",

	Long: useLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'Name'")
			cmd.Usage()
			os.Exit(1)
		}

		var name string

		if 0 < len(args) {

			name = args[0]

		}

		// Default body

		fmt.Println("hof config use", name)

	},
}
