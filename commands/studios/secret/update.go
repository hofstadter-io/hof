package secret

import (
	"fmt"

	// custom imports

	// infered imports
	"os"

	"github.com/hofstadter-io/hof/lib/secret"
	"github.com/spf13/cobra"
)

// Tool:   hof
// Name:   update
// Usage:  update <name> <env-file>
// Parent: secret

var UpdateLong = `Update a secret file that can be injected as environment variables`

var UpdateCmd = &cobra.Command{

	Use: "update <name> <env-file>",

	Short: "Update a secret",

	Long: UpdateLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In updateCmd", "args", args)
		// Argument Parsing
		// [0]name:   name
		//     help:
		//     req'd:  true
		if 0 >= len(args) {
			fmt.Println("missing required argument: 'name'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var name string

		if 0 < len(args) {

			name = args[0]
		}

		// [1]name:   file
		//     help:
		//     req'd:  true
		if 1 >= len(args) {
			fmt.Println("missing required argument: 'file'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var file string

		if 1 < len(args) {

			file = args[1]
		}

		/*
			fmt.Println("hof secret update:",
				name,

				file,
			)
		*/

		err := secret.Update(name, file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	// add sub-commands to this command when present

}
