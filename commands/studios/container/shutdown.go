package container

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
)

var shutdownLong = `Shutdown a Studios container.`

var ShutdownCmd = &cobra.Command{

	Use: "shutdown <name or id>",

	Short: "Shutdown a Studios container.",

	Long: shutdownLong,

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

		fmt.Println("hof studios container shutdown", ident)

	},
}
