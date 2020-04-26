package cmddatabase

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var statusLong = `Get the status of a Studios database.`

func StatusRun(ident string) (err error) {

	return err
}

var StatusCmd = &cobra.Command{

	Use: "status <name or id>",

	Short: "Get the status of a Studios database.",

	Long: statusLong,

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

		err = StatusRun(ident)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
