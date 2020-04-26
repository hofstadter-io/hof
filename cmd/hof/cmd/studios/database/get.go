package cmddatabase

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var getLong = `Get a Studios database`

func GetRun(ident string) (err error) {

	return err
}

var GetCmd = &cobra.Command{

	Use: "get <name or id>",

	Short: "Get a Studios database",

	Long: getLong,

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

		err = GetRun(ident)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
