package cmdcontainer

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var pushLong = `Push a Studios container.`

func PushRun(ident string) (err error) {

	return err
}

var PushCmd = &cobra.Command{

	Use: "push <name or id>",

	Short: "Push a Studios container.",

	Long: pushLong,

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

		err = PushRun(ident)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
