package secret

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hofstadter-io/hof/pkg/studios/secret"
)

var CreateLong = `Create a secret file that can be injected as environment variables`

var CreateCmd = &cobra.Command{

	Use: "create <name> <env-file>",

	Short: "Create a secret",

	Long: CreateLong,

	Run: func(cmd *cobra.Command, args []string) {

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'name'\n")
			cmd.Usage()
			os.Exit(1)
		}

		var name string
		if 0 < len(args) {
			name = args[0]
		}

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
