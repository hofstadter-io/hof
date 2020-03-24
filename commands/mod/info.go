package mod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"
)

var infoLong = `  print info about languages and modders known to mvs
    - no arg prints a list of known languages
    - an arg prints info about the language modder configuration that would be used`

var InfoCmd = &cobra.Command{

	Use: "info [language]",

	Short: "print info about languages and modders known to mvs",

	Long: infoLong,

	Run: func(cmd *cobra.Command, args []string) {

		// Argument Parsing

		var lang string

		if 0 < len(args) {

			lang = args[0]

		}

		msg, err := lib.LangInfo(lang)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(msg)

	},
}
