package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var getLong = `Get a Studios config`

var GetCmd = &cobra.Command{

	Use: "get <name or id>",

	Short: "Get a Studios config",

	Long: getLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'Ident'")
			cmd.Usage()
			os.Exit(1)
		}

		var ident string

		if 0 < len(args) {

			ident = args[0]

		}

		// Default body

		fmt.Println("hof studios config get", ident)

	},
}
