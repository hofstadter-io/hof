package mod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/mvs/lib"
)

var initLong = `initialize a new module in the current directory`

var InitCmd = &cobra.Command{

	Use: "init <lang> <module>",

	Short: "initialize a new module in the current directory",

	Long: initLong,

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
			fmt.Println("missing required argument: 'Module'")
			cmd.Usage()
			os.Exit(1)
		}

		var module string

		if 1 < len(args) {

			module = args[1]

		}

		err := lib.Init(lang, module)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}
