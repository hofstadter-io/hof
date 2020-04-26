package cmdapp

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var deleteLong = `Delete a Studios app.`

func DeleteRun(ident string) (err error) {

	return err
}

var DeleteCmd = &cobra.Command{

	Use: "delete <name or id>",

	Short: "Delete a Studios app.",

	Long: deleteLong,

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

		err = DeleteRun(ident)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
