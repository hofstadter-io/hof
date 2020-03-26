package mod

import (

	// hello... something might need to go here

	"os"

	"github.com/spf13/cobra"

	"fmt"

	"github.com/hofstadter-io/mvs/lib"
)

var convertLong = `convert another package system to MVS.`

var ConvertCmd = &cobra.Command{

	Use: "convert <lang> <file>",

	Short: "convert another package system to MVS.",

	Long: convertLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'Lang'")
			cmd.Usage()
			os.Exit(1)
		}

		var lang string

		if 0 < len(args) {

			lang = args[0]

		}

		if 1 >= len(args) {
			fmt.Println("missing required argument: 'Filename'")
			cmd.Usage()
			os.Exit(1)
		}

		var filename string

		if 1 < len(args) {

			filename = args[1]

		}

		err := lib.Convert(lang, filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
