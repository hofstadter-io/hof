package cmddatabase

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var saveLong = `Save a Studios database under a named reference.`

func SaveRun(ident string) (err error) {

	return err
}

var SaveCmd = &cobra.Command{

	Use: "save <name or id> <backup-name>",

	Short: "Save a Studios database under a named reference.",

	Long: saveLong,

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

		err = SaveRun(ident)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}
