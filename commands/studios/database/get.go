package database

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

var getLong = `Get a Studios database`

var GetCmd = &cobra.Command{

	Use: "get <name or id>",

	Short: "Get a Studios database",

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

		fmt.Println("hof studios database get", ident)

	},
}
