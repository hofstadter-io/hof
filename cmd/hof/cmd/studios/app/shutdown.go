package cmdapp

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var shutdownLong = `Shutdown a Studios app.`

func ShutdownRun(ident string) (err error) {

	return err
}

var ShutdownCmd = &cobra.Command{

	Use: "shutdown <name or id>",

	Short: "Shutdown a Studios app.",

	Long: shutdownLong,

	Run: func(cmd *cobra.Command, args []string) {
		var err error

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

		err = ShutdownRun(ident)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
