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
// Name:   create
// Usage:  create <name> <env-file>
// Parent: secret

var CreateLong = `Create a secret file that can be injected as environment variables`

var CreateCmd = &cobra.Command{

	Use: "create <name> <env-file>",

	Short: "Create a secret",

	Long: CreateLong,

	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("In createCmd", "args", args)
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
			fmt.Println("hof secret create:",
				name,

				file,
			)
		*/

		err := secret.Create(name, file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	},
}

func init() {
	// add sub-commands to this command when present

}
