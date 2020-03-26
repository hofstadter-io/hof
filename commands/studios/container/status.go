package container

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var statusLong = `Get the status of a Studios container.`

var StatusCmd = &cobra.Command{

	Use: "status <name or id>",

	Short: "Get the status of a Studios container.",

	Long: statusLong,

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

		fmt.Println("hof studios container status", ident)

	},
}
